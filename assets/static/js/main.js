
function copy(elt, id) {
  document.getElementById(id).select();
  document.execCommand("copy");
  elt.textContent = 'Copied!'
}
