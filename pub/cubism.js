/*
 * jank-flaccid cubism graph
 * route: /pub/cubism.html#min,max,step,size,hosts:hostsN,Filter
 */

$(document).ready(function() {

var app;
var hosts = [];
var opts = {};

var scripts = function(cb) {
	$('body').append('<script src="js/cubism.v1.js"></script>');
	cb();
}

var parse_hash = function() {
	var hash = window.location.hash.replace(/#/g,'');
	var arr = hash.split(',');
	if (arr.length < 6) {
		return false
	}

	opts.min = parseInt(arr[0],10);
	opts.max = parseInt(arr[1],10);
	opts.step = parseInt(arr[2],10);
	opts.size = parseInt(arr[3],10);
	opts.hosts = arr[4];
	opts.filter = arr[5];

	return true
}

var init = function() {

	truth = parse_hash();
	if (truth == false) {
		return
	}

	document.title = opts.filter

	app = new quickApp();
	app.Routes.Accounts_List(function(err,data) {
		console.log(data)
		hosts = data.Accounts;
		init_d3();
	});


}

var init_d3 = function() {
var context = cubism.context()
	.step(opts.step)
	.size(opts.size)

d3.select("body").selectAll(".axis")
    .data(["top", "bottom"])
  .enter().append("div")
    .attr("class", function(d) { return d + " axis"; })
    .each(function(d) { d3.select(this).call(context.axis().ticks(12).orient(d)); });

d3.select("body").append("div")
    .attr("class", "rule")
    .call(context.rule());

d3.select("body").selectAll(".horizon")
	.data(hosts.map(get_state3))

  .enter().insert("div", ".bottom")
    .attr("class", "horizon")
	.call(context.horizon().extent([opts.min, opts.max]));

context.on("focus", function(i) {
  d3.selectAll(".value").style("right", i == null ? null : context.size() - i + "px");
});


function pad(n) {
	return (n < 10 ? '0'+n : n)
}

function get_date(d) {
	var xd = new Date(d)
	return xd.getUTCFullYear() + "-" + (pad(xd.getUTCMonth()+1)) + "-" + (pad(xd.getUTCDate())) + "T" + (pad(xd.getUTCHours())) + ":" + (pad(xd.getUTCMinutes())) + ":" + (pad(xd.getUTCSeconds()))
}


function get_state3(x) {

	var value = 0,
	values = [],
	i = 0,
	last;

	return context.metric(function(start,stop,step,callback) {
		var dStart = get_date(start)
		var dStop = get_date(stop)
		d3.json("/query/state/"+x+"/nil/"+dStart+"/"+dStop+"/-1/"+opts.filter, function(data) {
			start = +start, stop = +stop;
			if (isNaN(last)) last = start;
			while (last < stop) {
				last += step;
//				values.push(value);
			}
			var res = []
			try {
				for(var v in data.Arr) {
					var datum = data.Arr[v]
					res.push(datum.Data)
				}
			} catch(err) {
			}
	callback(null, res);
	});
	}, x);
}


}


scripts(init);

});
