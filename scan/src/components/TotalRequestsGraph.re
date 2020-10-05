module Styles = {
  open Css;

  let card =
    style([
      backgroundColor(Colors.white),
      height(`percent(100.)),
      borderRadius(`px(4)),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(13)),
      Media.mobile([padding(`px(10))]),
    ]);

  let innerCard =
    style([width(`percent(100.)), height(`px(200)), margin2(~v=`zero, ~h=`auto)]);

  let infoHeader =
    style([
      borderBottom(`px(1), `solid, Colors.gray9),
      padding2(~h=`px(11), ~v=`zero),
      paddingBottom(`px(16)),
    ]);

  let chart = show => style([important(display(show ? `block : `none))]);
};

let renderGraph: array(HistoricalTotalRequestQuery.t) => unit = [%bs.raw
  {|
function(data) {
  var Chart = require('chart.js');
  var ctx = document.getElementById('historicalRequest').getContext('2d');

  // change seconds to milliseconds
  data = data.map(({y, t}) => {
    return {
      y: y,
      t: t * 1000,
    }
  });

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
        maintainAspectRatio: false,
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
              time: {
                unit: 'day',
                unitStepSize: 1,
                displayFormats: {
                  'day': 'MMM DD',
                },
              },
              ticks: {
                fontFamily: 'Inter',
                fontColor: '#888888',
                fontSize: 10,
                autoSkip: true,
                maxTicksLimit: 10,
              }
            },
          ],
          yAxes: [
            {
              gridLines: {
                display: true,
                color: "#f2f2f2",
                drawBorder: false,
                zeroLineColor: '#eaeaea'
              },
              ticks: {
                fontFamily: 'Inter',
                fontColor: '#888888',
                fontSize: 10,
                maxTicksLimit: 5,
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
  let dataQuery = HistoricalTotalRequestQuery.get();
  let (lastCount, setLastCount) = React.useState(_ => 0);

  React.useEffect1(
    () => {
      switch (dataQuery) {
      | Data(data) =>
        if (data->Belt.Array.size > 0) {
          // check the incoming data is a new data.
          let last = data->Belt.Array.get(data->Belt.Array.size - 1)->Belt.Option.getExn;
          if (last.y != lastCount) {
            setLastCount(_ => last.y);
            renderGraph(data);
          };
        }
      | _ => ()
      };
      None;
    },
    [|dataQuery|],
  );

  <div className=Styles.card>
    <div
      className={Css.merge([
        CssHelper.flexBox(),
        Styles.infoHeader,
        CssHelper.mb(~size=40, ()),
        CssHelper.mbSm(~size=16, ()),
      ])}>
      <Heading value="Total Requests" size=Heading.H4 />
      <HSpacing size=Spacing.xs />
      <CTooltip tooltipText="The total number of oracle data requests made">
        <Icon name="fal fa-info-circle" size=10 />
      </CTooltip>
    </div>
    {switch (dataQuery) {
     | Data(data) =>
       let show = data->Belt.Array.size > 5;
       <div className=Styles.innerCard>
         <canvas id="historicalRequest" className={Styles.chart(show)} />
         <EmptyContainer display={!show} height={`percent(100.)}>
           <Icon name="fal fa-clock" size=40 color=Colors.bandBlue />
           <VSpacing size={`px(16)} />
           <Heading
             size=Heading.H4
             value="Insufficient data to visualize"
             align=Heading.Center
             weight=Heading.Regular
             color=Colors.bandBlue
           />
         </EmptyContainer>
       </div>;
     | _ => <LoadingCensorBar fullWidth=true height=200 />
     }}
  </div>;
};
