window.onload = init;

function init() {
	rsMenu();
	currentPage();
}

function rsMenu() {
	/* Checking for DOM functionality */
	if (!document.getElementsByTagName('li')){ return; }
	if (!document.getElementsByTagName('nav')[0]){ alert('You forgot the NAV tag!');return; }
	if (!document.getElementsByTagName('ul')[0]){ alert('You forgot the UL tag!');return; }
	if (!document.createElement('div')) { return; }
	if (!document.createElement('img')) { return; }
	
	/* setting classes for the rs-menu */
	var jsNav = document.getElementsByTagName('nav')[0].className = 'js-nav';
	var jsUl = document.getElementsByTagName('ul')[0].className = 'js-ul';
	
	
	var nav = document.getElementsByTagName('nav')[0];
	var ul = document.getElementsByTagName('ul')[0];
	var linkItems = nav.getElementsByTagName('li');
	var links = nav.getElementsByTagName('a');

	/* Sets the Li elements to hidden, otherwise the show up immediatly when resizing */
	
	for(var i = 0; i < linkItems.length; i++) { 
  				linkItems[i].className = "hidden";
  	}
  	
	//Creates the menu div with the image 
	var menu = document.createElement("div");
	var img = document.createElement("img");
	img.id = "show"; 	img.src = "/images/menu-white.png";
	menu.id = "menu"; 	
	/* Adding the img to the menu div and the menu to the navigation. It will be inserted before the Ul */
	var para = document.createElement("p")
	var brand = document.createTextNode("Jorn Schalkwijk");
	para.appendChild(brand);
	menu.appendChild(img);
	menu.appendChild(para);
	nav.insertBefore(menu, ul)
	  	
  	/* Hidden is used to toggle the menu */
	var hidden = true;
	
	/* Responds to the click on the menu item and shows or hides the menu */
	show.addEventListener('click', function(event) {
		if(hidden == true) {
			for(var i = 0; i < linkItems.length; i++) {
			 	linkItems[i].className = "show";
  		    }
  			hidden = false;
  		} else {
  			for(var i = 0; i < linkItems.length; i++) { 
  				linkItems[i].className = "hidden";
  				}
  			hidden = true;	
  		}			
	});
}

function currentPage() {
	var path = window.location.pathname;
}

// functie verbeteren door de css ook door javascript te laten creëren.

