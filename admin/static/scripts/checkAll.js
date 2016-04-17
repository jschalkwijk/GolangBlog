/**
 * Created by Jorn on 15-4-2016.
 */

window.onload = init;

function init(){
    handler();
}

var c = document.getElementById("check-all");
var box = document.getElementsByName("checkbox");
var flag = false;
function handler(){
    c.addEventListener("click", function(e){
        checkAll();
    })
}
function checkAll(){
    if(!flag) {
        for (i = 0; i < box.length; i++) {
            box[i].checked = true;
            flag = true;
            c.innerHTML = "Un-Check"
        }
    } else if(flag){
        for (i = 0; i < box.length; i++) {
            box[i].checked = false;
            flag = false;
            c.innerHTML = "Check"
        }
    }
}
