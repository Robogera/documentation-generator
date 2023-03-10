const container = document.querySelector('#scrl');

let mouse_is_pressed;
let start_y;
let scroll_top;

container.addEventListener('mousedown', e => mousePressed(e));
container.addEventListener('mouseup', e => mouseReleased(e));
container.addEventListener('mouseleave', e => mouseLeftArea(e));
container.addEventListener('mousemove', e => mouseMoved(e));

function mousePressed(e){
  mouse_is_pressed = true;
  start_y = e.pageY - container.offsetTop;
  scroll_top = container.scrollTop;
}

function mouseReleased(e){
  mouse_is_pressed = false;
}

function mouseLeftArea(e){
  mouse_is_pressed = false;
}

function mouseMoved(e){
  if(mouse_is_pressed){
    e.preventDefault();
    const y = e.pageY - container.offsetTop;
    const walk_y = y - start_y;
    window.scrollBy({
      top: 100,
      behavior: "smooth",
    });
  }
}
