
function copy(id) {
  var copy = document.getElementById(id);
  copy.select();
  document.execCommand("Copy");
}
