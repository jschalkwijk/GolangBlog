/**
 * Created by jorn on 18-04-16.
 */
/* Made to add files on click to the tinyMCE editor */

window.onload = mce;

// called after image search.
function mce(){
    handleImagesAdding();
    check();
}

function insertImages(path,thumb) {
    if(thumb == null) {
        tinyMCE.execCommand('mceInsertRawHTML', false, '<a href='+ path +'><img src=' + path + ' width=100%></a>');
    } else {
        tinyMCE.execCommand('mceInsertRawHTML', false, '<a href='+ path +'><img src=' + thumb + '></a>');
    }
}

function handleImagesAdding(){

    if(document.getElementById("add-image")) {
        console.log("hello");
        document.getElementById("add-image").addEventListener("click", function (e) {
            // e.target is the clicked element!
            // If it was a list item
            if (e.target && e.target.className == "files") {
                // List item found!  Output the ID!
                insertImages(e.target.name);
            }
        });
    }

    if(document.getElementById("add-multiple")) {
        var string;
        var path;
        var thumb;
        document.getElementById("add-multiple").addEventListener("click", function (e) {
            var checked = document.getElementsByName("checkbox[]");
            // loop over them all
            for (var i = 0; i < checked.length; i++) {

                if (checked[i].checked) {
                    string = checked[i].value.split("#");
                    thumb = string[0];
                    path = string[1];
                    insertImages(path,thumb);
                }
            }
        });
    }
    if(document.getElementById("search-file")) {
        var search = document.getElementById("search-file");
        var change = document.getElementById("return");

        search.addEventListener("click", function () {
            var searchTerm = document.getElementById("search").value;

            xhttp = new XMLHttpRequest();
            xhttp.onreadystatechange = function () {
                change.innerHTML = xhttp.responseText;
                // after changing the inner html of the searh results, call the mce init function again
                mce();
            };

            xhttp.open("POST", "/admin/blocks/include-files-tinymce.php", true);
            xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
            xhttp.send("search-file=submit&search=" + searchTerm);
        });
    }
}
