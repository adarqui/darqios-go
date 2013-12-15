/*
 * jank-flaccid cubism graph
 * route: /pub/cubism.html#collection,min,max,step,size,hosts:hostsN,Filter,colors:colors
 *
 * example:
 *
 * http://ip:911/pub/cubism.html#state,0,20,5000,1000,nil,LoadAvg:last1min,ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:ffffff:00BBBB:00CCCC:00DDDD:00EEEE:00FFFF:00BB5E:00CC66:00DD6F:00EE77:00FF84:BBBB00:CCCC00:DDDD00:EEEE00:FFFF00J:BB5E00:CC6600m:DD6F00:EE7700:FF8000:BB0000:CC0000:DD0000:EE0000:FF0000
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

		opts.collection = arr[0];
		opts.min = parseInt(arr[1],10);
		opts.max = parseInt(arr[2],10);
		opts.step = parseInt(arr[3],10);
		opts.size = parseInt(arr[4],10);
		opts.hosts = arr[5];
		opts.filter = "b64="+btoa(arr[6]);
		opts.colors = arr[7];

		var colors = opts.colors.split(":")
		opts.colors = []
		for (var v in colors) {
		  var color = colors[v]
		  opts.colors.push('#'+color)
		}

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
			hosts = data.Accounts;
			init_d3();
		});
	}

	var init_d3 = function() {
		var context = cubism.context()
			.clientDelay(0)
			.serverDelay(0)
			.step(opts.step)
			.size(opts.size)

		d3.select("body").selectAll(".axis")
			.data(["top", "bottom"])
			.enter().append("div")
			.attr("class", function(d) { return d + " axis"; })
			.each(function(d) { d3.select(this).call(context.axis().ticks(10).orient(d)); });

		d3.select("body").selectAll(".horizon")
			.data(hosts.map(get_state3))
			.enter().insert("div", ".bottom")
			.attr("class", "horizon")
			.call(context.horizon().extent([opts.min, opts.max]).colors(opts.colors));

		context.on("focus", function(i) {
			d3.selectAll(".value").style("right", i == null ? null : window.innerWidth - i + "px");
		});


		function pad(n) {
			return (n < 10 ? '0'+n : n)
		}

		function get_date(d) {
			var xd = new Date(d)
			return xd.getUTCFullYear() + "-" + (pad(xd.getUTCMonth()+1)) + "-" + (pad(xd.getUTCDate())) + "T" + (pad(xd.getUTCHours())) + ":" + (pad(xd.getUTCMinutes())) + ":" + (pad(xd.getUTCSeconds()))
		}


		function get_state3(x) {

			return context.metric(function(start,stop,step,callback) {
				start = +start, stop = +stop;
				var dStart = get_date(start)
				var dStop = get_date(stop)
				d3.json("/query/"+opts.collection+"/"+x+"/nil/"+dStart+"/"+dStop+"/-1/"+opts.filter, function(data) {
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
