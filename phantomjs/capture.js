var page = require('webpage').create(), system = require('system');
var time, address, destination;

if (system.args.length != 3) {
    console.log('Usage: capture.js [some URL] [destination]');
    phantom.exit();
}

address = system.args[1];
destination = system.args[2];

page.open(address, function () {
    page.render(destination);
    phantom.exit();
});
