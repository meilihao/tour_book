# FAQ
## browser
### 空闲
```js
// https://gist.github.com/barraponto/4547ef5718fd2d31e5cdcafef0208096
const DOCUMENT_EVENTS = [
  'mousemove', 'mousedown', 'click',
  'touchmove', 'touchstart', 'touchend',
  'keydown', 'keypress'
];

export class IdleTimer {
  constructor(onIdleTimeout, timeout) {
    this.onIdleTimeout = onIdleTimeout;
    this.timeout = timeout;
    this.timer = null;
    this.active = false;
    this.resetTimer = this.resetTimer.bind(this);
  }

  activate() {
    if (!this.active) { this.bindEvents(); }
    this.timer = setTimeout(this.onIdleTimeout, this.timeout);
    this.active = true;
  }

  deactivate() {
    if (this.active) { this.unbindEvents(); }
    clearInterval(this.timer);
    this.active = false;
  }

  resetTimer() {
    clearInterval(this.timer);
    this.activate();
  }

  bindEvents() {
    window.addEventListener(
      'scroll', this.resetTimer, { capture: true, passive: true});
    window.addEventListener('load', this.resetTimer);
    DOCUMENT_EVENTS.forEach(
      eventType => document.addEventListener(eventType, this.resetTimer));
  }

  unbindEvents() {
    // remove only checks capture
    window.removeEventListener( 'scroll', this.resetTimer, { capture: true });
    window.removeEventListener('load', this.resetTimer);
    DOCUMENT_EVENTS.forEach(
      eventType => document.removeEventListener(eventType, this.resetTimer));
  }
}
```

使用:
```js
import { IdleTimer } from './idletimer';

const onTimeout = () => console.log('time is up!');
const idletimer = new IdleTimer(onTimeout, 7000);

idletimer.activate();
```