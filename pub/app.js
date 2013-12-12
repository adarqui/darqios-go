var quickApp = function(opts) {

	var APP = this;

	APP.opts = opts
	APP.settings = {
		ws : "filled_in",
		elm : {
			log : $("#log")
		},
	}

	APP.AppendLog = function(msg) {
		var d = APP.settings.elm.log[0]
		var doScroll = d.scrollTop == d.scrollHeight - d.clientHeight;
		msg.appendTo(APP.settings.elm.log)
		if (doScroll) {
			d.scrollTop = d.scrollHeight - d.clientHeight;
		}
	}

	APP.Routes = {
		Accounts_List : function(cb) {
			$.getJSON("/accounts:list", function(data,err) {
				cb(err,data);
			});
		},
		Query_Recent : function(x,cb) {
			$.getJSON("/query/state/"+x+"/nil/-6/now/-1/nil", function(data, err) {
				cb(err,data);
			});
		}
	}

	APP.WS_Init = function() {

		console.log("quickApp: WS_Init: Connecting");

		if (window["WebSocket"]) {
			APP.io = new WebSocket(APP.settings.ws);
			APP.io.onopen = function(evt) {
				if (APP.io_timeout != undefined) {
					clearTimeout(APP.io_timeout)
				}
				for(var v in APP.settings.channels) {
					var chan = APP.settings.channels[v]
					console.log("registering: ",chan)
					APP.io.send(JSON.stringify({ op: "register", channel : chan }))
				}
			}
			APP.io.onclose = function(evt) {
				APP.AppendLog($("<div><b>Connection closed.</b></div>"));
				APP.io_timeout = setTimeout(function() {
					console.log("Attempting to reconnect...")
					APP.io = new WebSocket(APP.settings.ws);
				}, 1000)
			}
			APP.io.onmessage = function(evt) {
				//APP.AppendLog($("<div/>").text(evt.data))
				try {
					jsn = JSON.parse(evt.data)
					var channel = jsn.Channel
				} catch(err) {
					return
				}
				for(var v in APP.io_callbacks) {
					var cb = APP.io_callbacks[v]
					cb(null,channel,evt.data)
				}
			}
		} else {
			APP.AppendLog($("<div><b>Your browser does not support WebSockets.</b></div>"))
		}
	}

	APP.Register = function(args) {
		if (typeof args === 'string') {
			APP.settings.channels[args] = args
		} else if (typeof args === 'object') {
			for (var v in args) {
				var arg = args[v]
				APP.settings.channels[arg] = arg
			}
		}
	}

	APP.Callback = function(cb) {
		APP.io_callbacks.push(cb)
	}

	APP.Init = function() {

	   APP.settings.channels = {}
	   APP.io_callbacks = []
	   var url = (window.location.origin)
	   APP.settings.ws = (url+"/ws").replace(/http/g,"ws")
	   APP.settings.url = url

	   console.log("quickApp: Initialized")
	}


	APP.GrabSingle = function(data, path) {
		var nd = {}
		try {
			var split_path = path.split('.')
			var new_path = []
			for (var v in split_path) {
				new_path.push('["'+split_path[v]+'"]')
			}
			new_path = new_path.join('')
			eval("var _o = data"+new_path+";")
//				if (_o == 0 ) return nd;
				nd = { host: data.Host, data: _o }
			} catch(err) {
		}
		return nd
	}

	APP.GrabArray = function(data, path) {
		var nd = {}
		for (var v in data) {
			var datum = data[v]
			try {
				eval("var _o = datum."+path+";")
//				if ( _o == 0 ) { continue; }
				nd[datum.Host] = { host: datum.Host , data: _o }
			} catch(err) {
				continue
			}
		}
		return nd
	}

	APP.StateGrab = function(accounts, path, cb) {
		$.getJSON('/accounts/get/'+accounts, function(data,err) {
			/* FUCKOFF */
			var nd = {}
			for (var v in data) {
				var datum = data[v]
				try {
					eval("var _o = datum."+path+";")
					if ( _o == 0 ) { continue; }
					nd[datum.Host] = { host: datum.Host , data: _o }
				} catch(err) {
				}
			}
			var nd = APP.GrabArray(data, path)
			cb(err,nd)
		})
	}


	APP.State2CSV = function(type,name,cb) {
		$.get('/accounts/get/'+name, function(data,err) {

			data = JSON.parse(data);
			var csv = []
			var i = 0
			for (var v in data) {
				var account = data[v]
				try {

					var elm
					if(type == "load") {
						elm = { name : account.Host, value : account.State.LoadAvg.last1min };
					}
					else if(type == "procTot") {
						elm = { name : account.Host, value : account.State.Proc.Total };
					} else {
						return
					}
					csv[i] = elm
					i = i + 1;
				} catch (err) { }
			}

			data = csv
			cb(err,data)
		});
	}

	APP.D3_Vertical_Bar = function(hw,data) {
		var w = hw[0]
		var h = hw[1]
		var barpad = 1

		var svg

		if (APP.d3verticalbar == undefined) {

			svg = d3.select("body")
				.append("svg")
				.attr("width",w)
				.attr("height",h)

			APP.d3verticalbar = svg

		} else {
			svg = APP.d3verticalbar

			svg.selectAll("rect").data(APP.olddata).remove()
			/*
			svg.selectAll("rect")
			.data(data)
			.transition()
			.duration(500)
			*/

		}

		svg.selectAll("rect")
		.data(data)
		.enter()
		.append("rect")
		.attr("x", function(d, i) {
			return i * (w / data.length)
		})
		.attr("y", function(d) { return h - d.value; })
		.attr("width", w / data.length - barpad)
		.attr("height", function(d) { return d.value; })
		.attr("fill", function(d) {
			return "rgb(0,0,"+d.value+")";
		});

		svg.selectAll("text")
		.data(data)
		.enter()
		.append("text")
		.text(function(d) {
			return d.name
		})
		/*
		.attr("x", function(d, i) {
			return i * (w / data.length)
		})
		.attr("y", function(d) {
			return h + 20
		})
		*/
	   .attr("x", function(d,i) {
		   return i * (w / data.length)
	   })
	   .attr("y", 10)
		.attr("font-family", "sans-serif")
		.attr("font-size", "11px")
		.attr("fill","red");

//		svg.selectAll("rect").data(data).exit()

		APP.olddata = data
	}

	APP.D3_Horizontal_Bar = function(data) {

		console.log("D3_Horizontal_Bar")

var m = [40, 10, 10, 40],
    w = 960 - m[1] - m[3],
    h = 930 - m[0] - m[2];

var format = d3.format(",.0f");

var x = d3.scale.linear().range([0, w]),
    y = d3.scale.ordinal().rangeRoundBands([0, h], .1);

var xAxis = d3.svg.axis().scale(x).orient("top").tickSize(-h),
    yAxis = d3.svg.axis().scale(y).orient("left").tickSize(0);

var svg = d3.select("body").append("svg")
    .attr("width", w + m[1] + m[3])
    .attr("height", h + m[0] + m[2])
  .append("g")
    .attr("transform", "translate(" + m[3] + "," + m[0] + ")");

  // Parse numbers, and sort by value.
  data.forEach(function(d) { d.value = +d.value; });
  data.sort(function(a, b) { return b.value - a.value; });

  // Set the scale domain.
  x.domain([0, d3.max(data, function(d) { return d.value; })]);
  y.domain(data.map(function(d) { return d.name; }));

  var bar = svg.selectAll("g.bar")
      .data(data)
    .enter().append("g")
      .attr("class", "bar")
      .attr("transform", function(d) { return "translate(0," + y(d.name) + ")"; });

  bar.append("rect")
      .attr("width", function(d) { return x(d.value); })
      .attr("height", y.rangeBand());

  bar.append("text")
      .attr("class", "value")
      .attr("x", function(d) { return x(d.value); })
      .attr("y", y.rangeBand() / 2)
      .attr("dx", -3)
      .attr("dy", ".35em")
      .attr("text-anchor", "end")
      .text(function(d) { return format(d.value); });

  svg.append("g")
      .attr("class", "x axis")
      .call(xAxis);

  svg.append("g")
      .attr("class", "y axis")
      .call(yAxis);
}



APP.D3Line = function(hw,data) {

	var w = hw[0]
	var h = hw[1]
	var barpad = 1

	var svg

	if (APP.d3linesvg == undefined) {

var margin = {top: 20, right: 20, bottom: 30, left: 50},
    width = 960 - margin.left - margin.right,
    height = 500 - margin.top - margin.bottom;

var parseDate = d3.time.format("%d-%b-%y").parse;

var x = d3.time.scale()
    .range([0, width]);

var y = d3.scale.linear()
    .range([height, 0]);

var xAxis = d3.svg.axis()
    .scale(x)
    .orient("bottom");

var yAxis = d3.svg.axis()
    .scale(y)
    .orient("left");

var line = d3.svg.line()
    .x(function(d) { console.log("x",d); return x(d.name); })
    .y(function(d) { console.log("y",d); return y(d.value); });

var svg = d3.select("body").append("svg")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
  .append("g")
    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

	APP.d3linesvg = svg
	}
	else {
		svg = APP.d3linesvg
		svg.selectAll("body").data(APP.olddata).remove()
	}

   /*
$.get('/accounts/get/certs', function(data) {


	data = JSON.parse(data);
	var csv = []
	var i = data.length-1

	var state = data[0].State

	csv[0] = { date: 0, close: 100}
	csv[1] = { date: 1, close: 200}
	csv[2] = { date: 2, close: 150}
	data = csv
	console.log(csv)
	*/

  x.domain(d3.extent(data, function(d) { return d.name; }));
  y.domain(d3.extent(data, function(d) { return d.value; }));

  svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0," + height + ")")
      .call(xAxis);

  svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)
    .append("text")
      .attr("transform", "rotate(-90)")
      .attr("y", 6)
      .attr("dy", ".71em")
      .style("text-anchor", "end")
      .text("Price ($)");

  svg.append("path")
  		.datum(data)
//      .data(data)
      .attr("class", "line")
      .attr("d", line);

	  APP.olddata = data

}

	APP.Init()
}
