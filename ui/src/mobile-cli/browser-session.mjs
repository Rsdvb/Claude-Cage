import { browserRequest, desktopRequest } from "./transport.mjs";

const clamp = (value, low, high) => Math.max(low, Math.min(high, value));
export function viewportSize(element) {
  return {
    width: clamp(Math.round(element.clientWidth), 240, 1440),
    height: clamp(Math.round(element.clientHeight), 160, 1600),
  };
}

export function framePoint(rect, frame, x, y) {
  const width = frame.content_width || frame.width;
  const height = frame.content_height || frame.height;
  const scale = Math.min(rect.width / width, rect.height / height);
  if (!(scale > 0)) return null;
  const left = rect.left + (rect.width - width * scale) / 2;
  const top = rect.top + (rect.height - height * scale) / 2;
  const px = Math.floor((x - left) / scale);
  const py = Math.floor((y - top) / scale);
  return px >= 0 && px < width && py >= 0 && py < height ? {x: px, y: py, scale} : null;
}

// Only frame reads repeat. One input or resize is followed by a fresh frame
// before the next input is accepted. A failed write is never replayed.
export function connectDisplay(stage, update, id, start, tool, initialURL) {
  const desktop = tool === "desktop";
  const request = desktop ? desktopRequest : browserRequest;
  const toolName = desktop ? "Claude Desktop" : "远程浏览器";
  let waitingSince = 0, homeNotice = "";
  const base = "/" + id;
  let alive = true, visible = true, opened = false, opening = true, writing = false;
  let frame = null, read, timer, paused = false, mayExist = false;
  let generation = 0, actionFrame = false, quietReads = 0, writeType = "";
  let pixelRatioSupported = false;
  const displaySize = () => ({...viewportSize(stage),
    ...(pixelRatioSupported ? {pixel_ratio:clamp(window.devicePixelRatio || 1, 1, 2)} : {})});
  const needsResize = size => frame && (frame.width !== size.width || frame.height !== size.height ||
    (pixelRatioSupported && frame.requested_pixel_ratio !== size.pixel_ratio));
  let queuedScroll = null, scrollTimer;
  // Only untransmitted displacement lives in this queue. A scroll is taken
  // out before sending; neither an interrupted write nor an old document replays it.
  function cancelScroll() {
    clearTimeout(scrollTimer); scrollTimer = null; queuedScroll = null;
  }
  let awaitingResizeFrame = false, pendingAddress = null;
  const lifetime = new AbortController();
  function discardAddress() {
    awaitingResizeFrame = false;
    pendingAddress?.resolve(false);
    pendingAddress = null;
  }
  const closeRemote = () => {
    if (mayExist) void request(base + "/close", {body: {}, keepalive: true}).catch(() => {});
  };
  const publish = (state) => { if (alive && visible) update(homeNotice && state.frame ? {...state, homeNotice} : state); };
  const cancelRead = () => { clearTimeout(timer); read?.abort(); read = null; };
  function failure(error, operation) {
    discardAddress(); cancelScroll();
    frame = null;
    paused = true;
    if (!desktop && operation === "frame" && error.status === 503) opened = mayExist = false;
    const reopen = !opened || [401, 404, 410, 423].includes(error.status) || (error.status === 403 && operation !== "action");
    const message = error.status === 404 && !opened ? "这台电脑尚不支持" + toolName + "，请更新电脑端 Cage。" :
      error.message;
    publish({error: message, uncertain: operation === "action" && !error.status, reopen});
  }
  const sameSurface = (first, second) => first && second && first.document === second.document &&
    first.width === second.width && first.height === second.height &&
    first.content_width === second.content_width && first.content_height === second.content_height;
  function finishFrame(changed, fromAction = false) {
    awaitingResizeFrame = false;
    if (pendingAddress) {
      const {url, resolve} = pendingAddress;
      pendingAddress = null;
      void send({type:"navigate", url}, frame).then(resolve);
      return;
    }
    if (changed || fromAction) publish({frame, busy:false});
    const size = displaySize();
    if (needsResize(size)) {
      cancelScroll();
      timer = setTimeout(resize, 0);
      return;
    }
    if (queuedScroll && !sameSurface(queuedScroll.expected, frame)) cancelScroll();
    if (queuedScroll) { flushScroll(); return; }
    quietReads = changed || fromAction ? 0 : Math.min(5, quietReads + 1);
    // The first follow-up catches delayed paint after an input. Stable pages
    // use conditional reads and progressively return to the idle cadence.
    const delay = (desktop ? [450, 700, 1100, 1600, 2200, 2500] : [90, 180, 380, 800, 1500, 2500])[quietReads];
    timer = setTimeout(readFrame, delay);
  }
  async function decodeFrame(next, valid) {
    if (next?.id !== id || typeof next.image !== "string" || !next.image ||
        typeof next.document !== "string" || !Number.isInteger(next.width) || next.width <= 0 ||
        !Number.isInteger(next.height) || next.height <= 0) throw new Error("主机返回的画面无效。");
    if ((next.content_width !== undefined || next.content_height !== undefined) &&
        (!Number.isFinite(next.content_width) || next.content_width <= 0 ||
         !Number.isFinite(next.content_height) || next.content_height <= 0)) throw new Error("主机返回的画面尺寸无效。");
    const src = "data:image/jpeg;base64," + next.image;
    const changed = frame?.image !== src || !sameSurface(frame, next) || frame.focus !== next.focus ||
      frame.url !== next.url || frame.title !== next.title || frame.password !== next.password || frame.window !== next.window || frame.shared !== next.shared ||
      frame.login_supported !== next.login_supported || frame.content_width !== next.content_width || frame.content_height !== next.content_height;
    if (frame?.image !== src) {
      const image = new Image();
      image.src = src;
      await image.decode();
    }
    if (!valid()) return null;
    if (next.url && next.url !== "about:blank") homeNotice = "";
    if (!desktop && next.action_frame === true) actionFrame = true;
    pixelRatioSupported = desktop && next.pixel_ratio_supported === true;
    waitingSince = 0;
    frame = {...next, image:src};
    return changed;
  }
  async function readFrame() {
    if (!alive || !visible || !opened || paused || writing || read) return;
    const controller = new AbortController();
    read = controller;
    const valid = () => read === controller && !controller.signal.aborted && alive && visible && !paused;
    const known = frame;
    const tag = actionFrame && /^[a-f0-9]{32}$/.test(known?.frame_tag || "") ? known.frame_tag : "";
    try {
      const next = await request(base + "/frame" + (tag ? "?after=" + tag : ""), {signal:controller.signal});
      if (!valid()) return;
      let changed = false;
      if (next?.unchanged === true) {
        if (!tag || frame !== known || next.id !== id || next.frame_tag !== tag) throw new Error("远程画面的版本不匹配，请重新读取。");
      } else {
        changed = await decodeFrame(next, valid);
        if (changed === null) return;
      }
      read = null;
      finishFrame(changed);
    } catch (error) {
      if (valid()) {
        if (error.status === 425 && (!waitingSince || Date.now() - waitingSince < 90000)) {
          waitingSince ||= Date.now();
          frame = null;
          publish({waiting:true});
          timer = setTimeout(readFrame, 1000);
        } else {
          failure(error, "frame");
        }
      }
    } finally {
      if (read === controller) read = null;
    }
  }
  function canSend(expected, focus) {
    return alive && visible && !writing && !paused && frame && expected &&
      sameSurface(frame, expected) &&
      (!focus || focus === frame.focus);
  }
  async function send(payload, expected, operation = "action") {
    if (!canSend(expected, payload.focus)) return false;
    if (["navigate","back","forward","reload"].includes(payload.type)) homeNotice = "";
    if (payload.type !== "scroll") cancelScroll();
    cancelRead();
    const current = generation;
    const valid = () => alive && visible && !paused && current === generation;
    writing = true; writeType = payload.type;
    if (operation === "resize") awaitingResizeFrame = true;
    publish({frame, busy:true});
    let feedback = false;
    try {
      const body = {...payload, document:expected.document};
      if (actionFrame) body.return_frame = true;
      const result = await request(base + "/action", {body, signal:lifetime.signal});
      if (!valid()) return false;
      if (actionFrame && result?.image) {
        feedback = await decodeFrame(result, valid) !== null;
      } else {
        if (result?.id !== id) throw new Error("主机未确认这次操作，请重新读取画面。");
        frame = null;
      }
      return valid();
    } catch (error) {
      if (valid()) failure(error, operation);
      return false;
    } finally {
      writing = false; writeType = "";
      if (alive && visible && !paused) {
        if (feedback && valid()) finishFrame(true, true);
        else void readFrame();
      }
    }
  }
  function flushScroll() {
    clearTimeout(scrollTimer); scrollTimer = null;
    if (!queuedScroll || writing || !frame || paused || !alive || !visible) return;
    const next = queuedScroll;
    if (!sameSurface(next.expected, frame)) { cancelScroll(); return; }
    const dx = clamp(next.dx, -3200, 3200), dy = clamp(next.dy, -3200, 3200);
    next.dx -= dx; next.dy -= dy;
    if (!next.dx && !next.dy) queuedScroll = null;
    void send({type:"scroll", x:next.x, y:next.y, dx, dy}, frame);
  }
  function scroll(expected, {x,y,dx,dy}) {
    if (desktop || !alive || !visible || paused || !opened || (!dx && !dy) ||
        (writing && writeType !== "scroll") || (frame && !sameSurface(expected, frame))) return;
    if (!queuedScroll || queuedScroll.x !== x || queuedScroll.y !== y || !sameSurface(queuedScroll.expected, expected)) {
      queuedScroll = {expected, x, y, dx:0, dy:0};
    }
    queuedScroll.dx += dx; queuedScroll.dy += dy;
    if (!queuedScroll.dx && !queuedScroll.dy) { cancelScroll(); return; }
    // Coalesce moves within one display tick. Ongoing writes drain the still
    // unsent movement once their next authenticated frame is available.
    if (!scrollTimer) scrollTimer = setTimeout(() => {scrollTimer = null; flushScroll();}, 16);
  }
  function resize() {
    if (!frame || !alive || !visible || writing || paused) return;
    const size = displaySize();
    if (needsResize(size)) {
      void send({type: "resize", ...size}, frame, "resize");
    }
  }
  async function open() {
    try {
      const body = {id, ...viewportSize(stage)};
      if (tool === "browser" && initialURL) body.url = initialURL;
      if (tool === "claude") body.url = "https://claude.ai/";
      const result = await request("/open", {body});
      if (result?.id !== id) throw new Error("主机未确认这次画面连接。");
      if (tool === "browser" && typeof result.home_notice === "string") homeNotice = result.home_notice;
      actionFrame = !desktop && result.action_frame === true;
      opened = true; mayExist = true;
      void readFrame();
    } catch (error) {
      mayExist = !error.status || error.status >= 500;
      if (alive) failure(error, "open");
    } finally {
      opening = false;
      if (!alive) void closeRemote();
    }
  }
  const observer = new ResizeObserver(() => {
    // Rotation/keyboard animations can report several intermediate sizes.
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(resize, 160);
  });
  let resizeTimer;
  observer.observe(stage);
  if (start) {
    void open();
  } else {
    // History navigation or a page reload resumes the known session by read.
    // A new host browser always requires another explicit open gesture.
    opened = true; opening = false; mayExist = true;
    void readFrame();
  }
  return {
    send, canSend, scroll, cancelScroll,
    async navigate(url) {
      if (!alive || !visible || !opened || paused) throw new Error("请先恢复远程画面，再前往网址。");
      if (awaitingResizeFrame && !pendingAddress) {
        return new Promise(resolve => {pendingAddress = {url, resolve};});
      }
      if (writing || !frame || pendingAddress) throw new Error("上一项操作尚未完成，请稍后再次前往。");
      return send({type: "navigate", url}, frame);
    },
    offline() {
      if (!alive) return;
      generation++; cancelRead(); clearTimeout(resizeTimer);
      failure(new Error("网络已断开，恢复后请重新读取画面。"), writing && writeType !== "resize" ? "action" : "frame");
    },
    refresh() {
      if (!opened || !alive || !visible) return;
      paused = false; waitingSince = 0;
      publish({});
      void readFrame();
    },
    visibility(active) {
      if (active === visible) return;
      visible = active;
      if (!visible) {
        generation++; cancelScroll(); discardAddress();
        cancelRead(); clearTimeout(resizeTimer); frame = null;
        if (alive) update({});
      } else if (alive) {
        if (!opened && !opening) publish({error: "远程画面未连接。", reopen: true});
        else if (paused) publish({error: "请重新读取远程画面。", reopen: !opened});
        else void readFrame();
      }
    },
    close() {
      if (!alive) return;
      alive = false; generation++; cancelScroll();
      discardAddress();
      cancelRead(); clearTimeout(resizeTimer); observer.disconnect();
      frame = null;
      lifetime.abort();
      // If opening is still in flight, its finally sends close after the
      // server has processed open. An early close could miss the new session.
      if (!opening) void closeRemote();
    },
  };
}
