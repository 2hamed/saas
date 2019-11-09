var page = require('webpage').create(), system = require('system'), sha1 = require('./sha1.js');
var time, address;

if (system.args.length === 1) {
    console.log('Usage: capture.js [some URL]');
    phantom.exit();
}

time = Date.now();
address = system.args[1];

var hash = sha1.create();
hash.update(address)

fileName = hash.hex() + "_" + time;

page.open(address, function () {
    page.render(fileName + ".png");
    phantom.exit();
});
