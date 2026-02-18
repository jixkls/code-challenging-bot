// Custom cursor (desktop only)
const isMobile = window.innerWidth <= 640;
const cursor = document.getElementById('cursor');
const ring = document.getElementById('cursorRing');

if (!isMobile) {
  let mx=0,my=0,rx=0,ry=0;
  document.addEventListener('mousemove',e=>{
    mx=e.clientX; my=e.clientY;
    cursor.style.left=mx+'px'; cursor.style.top=my+'px';
  });
  function animRing(){
    rx+=(mx-rx)*0.12; ry+=(my-ry)*0.12;
    ring.style.left=rx+'px'; ring.style.top=ry+'px';
    requestAnimationFrame(animRing);
  }
  animRing();
  document.querySelectorAll('a,button,.feature-card,.lang-card,.level-card,.step,.phase-dot,.cmd-item,.diff-tag').forEach(el=>{
    el.addEventListener('mouseenter',()=>ring.classList.add('hovering'));
    el.addEventListener('mouseleave',()=>ring.classList.remove('hovering'));
  });
}

// Intersection Observer for reveals
const observer = new IntersectionObserver(entries=>{
  entries.forEach(e=>{
    if(e.isIntersecting) { e.target.classList.add('visible'); }
  });
},{threshold:0.12});
document.querySelectorAll('.reveal').forEach(el=>observer.observe(el));

// Nav shrink on scroll
window.addEventListener('scroll',()=>{
  document.querySelector('nav').style.padding = window.scrollY>60 ? '12px 0' : '18px 0';
});

// Mobile menu
function toggleMenu() {
  const links = document.getElementById('navLinks');
  const hamburger = document.getElementById('navHamburger');
  links.classList.toggle('open');
  hamburger.classList.toggle('active');
  document.body.style.overflow = links.classList.contains('open') ? 'hidden' : '';
}
function closeMenu() {
  const links = document.getElementById('navLinks');
  const hamburger = document.getElementById('navHamburger');
  links.classList.remove('open');
  hamburger.classList.remove('active');
  document.body.style.overflow = '';
}
