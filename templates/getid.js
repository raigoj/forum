const posts = document.querySelector('#post-id');
posts.dataset.postId
fetch('/posts')
    .then(response => response.text())
    .then(data => console.log(data));

