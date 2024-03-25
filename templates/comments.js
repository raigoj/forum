const postId = document.getElementById('post-id');
postId.addEventListener('click', event => {
console.log(event.target.getAttribute('data-id'));

});

