module Styles = {
  open Css;

  let card =
    style([
      backgroundColor(Colors.white),
      height(`percent(100.)),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      Media.mobile([padding(`px(10))]),
    ]);

  let innerCard =
    style([width(`percent(100.)), maxWidth(`px(420)), margin2(~v=`zero, ~h=`auto)]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);
};

let renderGraph: unit => unit = [%bs.raw
  {|
function() {
  var Chart = require('chart.js');
  var ctx = document.getElementById('historicalRequest').getContext('2d');

  // mock data
  let data = [];
  let x = 0;
  for (let i = 100; i > 0; i--) {
    x += Math.floor(Math.random() * 10000);
    data.push({t: Date.now() - (i * 86400000), y: x});
  }

  var chart = new Chart(ctx, {
      // The type of chart we want to create
      type: 'line',

      // The data for our dataset
      data: {
          datasets: [{
              type: 'line',
              pointRadius: 0,
              fill: false,
              borderColor: '#5269ff',
              data: data,
              borderWidth: 2,
          }]
      },

      // Configuration options go here
      options: {
        legend: {
          display: false,
        },
        scales: {
          xAxes: [
            {
              type: 'time',
              distribution: 'series',
              gridLines: {
                display: false,
                drawBorder: false,
              },
              ticks: {
                fontFamily: 'Inter',
                fontColor: '#888888',
                fontSize: 10,
                autoSkip: true,
                maxTicksLimit: 15,
              }
            },
          ],
          yAxes: [
            {
              ticks: {
                fontFamily: 'Inter',
                fontColor: '#888888',
                fontSize: 10,
                stepSize: 100000,
                callback: function(value) {
                  var ranges = [
                      { divider: 1e6, suffix: 'M' },
                      { divider: 1e3, suffix: 'K' }
                  ];
                  function formatNumber(n) {
                      for (var i = 0; i < ranges.length; i++) {
                        if (n >= ranges[i].divider) {
                            return (n / ranges[i].divider).toString() + ranges[i].suffix;
                        }
                      }
                      return n;
                  }
                  return formatNumber(value);
                }
              }
            },
          ],
        },
        tooltips: {
					mode: 'index',
					intersect: false,
          backgroundColor: '#555555',
          titleFontFamily: "Inter",
          titleFontSize: 12,
          titleFontColor: '#ffffff',
          titleFontStyle: "500",
          titleMarginBottom: 2,
          bodyFontFamily: "Inter",
          bodyFontSize: 10,
          bodyFontColor: '#888888',
          bodyFontStyle: "normal",
          xPadding: 15,
          yPadding: 10,
          caretSize: 6,
          displayColors: false,
          callbacks: {
            title: function(tooltipItem, data) {
              var title = (parseInt(tooltipItem[0].value)).toLocaleString();
              return title + " requests";
            },
            label: function(tooltipItem, data) {
              let date = new Date(tooltipItem.label);
              let dateTimeFormat = new Intl.DateTimeFormat('en', { year: 'numeric', month: 'short', day: '2-digit' });
              let [{ value: month },,{ value: day },,{ value: year }] = dateTimeFormat .formatToParts(date );

              return `${month} ${day},${year}`;
            },
          }
				},
      }
  });
}
  |}
];

[@react.component]
let make = () => {
  React.useEffect0(() => {
    renderGraph();
    None;
  });

  <div className=Styles.card>
    <Heading value="Information" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
    <div className=Styles.innerCard> <canvas id="historicalRequest" /> </div>
  </div>;
};
