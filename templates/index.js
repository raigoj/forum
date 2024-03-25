const btn = document.getElementById('btn');

btn.addEventListener('click', event => {
  // 👇️ {bar: 'foo'}
  console.log(event.target.dataset);

  // 👇️ "foo"
  console.log(event.target.getAttribute('data-bar'));
});