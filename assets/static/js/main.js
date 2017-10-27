
function copy(elt, id) {
  document.getElementById(id).select();
  document.execCommand("Copy");
  elt.textContent = 'Copied!'
}
