var context = cubism.context()
    .step(1e4)
    .size(1440);

d3.select("body").selectAll(".axis")
    .data(["top", "bottom"])
  .enter().append("div")
    .attr("class", function(d) { return d + " axis"; })
    .each(function(d) { d3.select(this).call(context.axis().ticks(12).orient(d)); });

d3.select("body").append("div")
    .attr("class", "rule")
    .call(context.rule());

d3.select("body").selectAll(".horizon")
    .data(d3.range(1, 50).map(random))
  .enter().insert("div", ".bottom")
    .attr("class", "horizon")
    .call(context.horizon().extent([-10, 10]));

context.on("focus", function(i) {
  d3.selectAll(".value").style("right", i == null ? null : context.size() - i + "px");
});

// Replace this with context.graphite and graphite.metric!
function random(x) {
  var value = 0,
      values = [],
      i = 0,
      last;
  return context.metric(function(start, stop, step, callback) {
    start = +start, stop = +stop;
    if (isNaN(last)) last = start;
    while (last < stop) {
      last += step;
      value = Math.max(-10, Math.min(10, value + .8 * Math.random() - .4 + .2 * Math.cos(i += x * .02)));
	  value = 1
      values.push(value);
    }
	console.log("Values", values, "start", start, "stop", stop)
    callback(null, values = values.slice((start - stop) / step));
  }, x);
}
