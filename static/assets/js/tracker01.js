var slider = document.getElementById("apparition");
var output = document.getElementById("demo");
output.innerHTML = slider.value;

slider.oninput = function() {
    output.innerHTML = this.value;
}

var sliderr = document.getElementById("album");
var outputt = document.getElementById("demoo");
outputt.innerHTML = sliderr.value;

sliderr.oninput = function() {
    outputt.innerHTML = this.value;
}