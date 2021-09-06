
var cookies = {}
var pins = process.env.pins
if(pins){
	pins = pins.split("&")
	for (var key in cookies) {
	    c = false
	    for (var pin of pins) {
		   if (pin && cookies[key].indexOf(pin) != -1) {
			  c = true
			  break
		   }
	    }
	    if (!c) {
		   delete cookies[key]
	    }
	}
}
module.exports = cookies