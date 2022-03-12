openorclose = 0
Form()
function Form() {
    if (openorclose == 1) {
        document.getElementById("myForm").style.display = "block";
        openorclose = 0
    } else {
        document.getElementById("myForm").style.display = "none";
        openorclose = 1
    }
}

