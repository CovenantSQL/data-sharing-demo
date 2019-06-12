<template>
  <div id="arch">
    <div id="chart"></div>
  </div>
</template>

<script>
  export default {
    name: "Arch",
    data() {
      return {
        frame_num: 0,
        currentFrame: null,
      }
    },
    mounted() {
    },
    methods: {
      load: function (cql) {
        if (this.frame_num === 0) {
          window.archPlayer.frame(this.frame_num, "Log Replication" + this.frame_num, this.initFrame());
        }
        // window.archPlayer.current().model().clients.create("Y" + this.frame_num);
        // window.archPlayer.frame(this.frame_num, "Log Replication" + this.frame_num, this.newFrame());
        window.archPlayer.next()
        this.newFrame(cql)(window.archPlayer.current())
        this.frame_num += 1;
      },
      initFrame: function () {
        return function (fr) {
          this.currentFrame = fr;
          var
            player = fr.player(),
            layout = fr.layout(),
            model = function () {
              return fr.model();
            },
            cluster = function (value) {
              model().nodes.toArray().forEach(function (node) {
                node.cluster(value);
              });
            };

          fr
            .after(0, function () {
              model().clear();
            })
            .after(0, function () {
              fr.model().title = '<h2 style="visibility:visible">Log Replication</h1>'
                + '<br/>' + fr.model().controls.html();
            })
            .after(0, function () {
              model().title = "";
            })

            //------------------------------
            // Cluster Initialization
            //------------------------------
            .after(0, function () {
              model().nodes.create("A");
              model().nodes.create("B");
              model().nodes.create("C");
              cluster(["A", "B", "C"]);
              model().clients.create("X");
              layout.invalidate();
            })
            .after(200, function () {
              model().forceImmediateLeader();
            })

          //------------------------------
          // Single Entry Replication
          //------------------------------
          // .after(300, function () {
          //   client("X").send(model().leader(), "select * from cargos where id = 5");
          // })
          // .at(model(), "appendEntriesRequestsSent", function () {
          // })
          // .at(model(), "commitIndexChange", function (event) {
          //   if (event.target !== model().leader()) {
          //     subtitle('<h2>commitIndexChange</h2>');
          //   }
          // });
          // .after(model().defaultNetworkLatency * 2, function (event) {
          //   subtitle('<h2>finish</h2>');
          // })

          player.play();
        }
      },
      newFrame: function (cql) {
        return function (fr) {
          const player = fr.player(),
            layout = fr.layout(),
            model = function () {
              return fr.model();
            },
            client = function (id) {
              return fr.model().clients.find(id);
            },
            subtitle = function (s, pause) {
              model().subtitle = s + model().controls.html();
              layout.invalidate();
              if (pause === undefined) {
                model().controls.show()
              }
            };

          fr.after(500, function () {
            subtitle('<h2>' + cql +'</h2>', false);
            client("X").send(model().leader(), cql);
          })
          .at(model(), "appendEntriesRequestsSent", function () {
          })
          .at(model(), "commitIndexChange", function (event) {
            return (event.target !== model().leader())
          })
          .after(model().defaultNetworkLatency * 2, function () {
            subtitle('<h2>Committed</h2>');
          })

          player.play();
        }
      }
    }
  }
</script>

<style scoped>
  svg rect {
    shape-rendering: crispEdges
  }

  #chart {
    min-height: 500px;
    height: 500px;
  }

  .log-entry {
    text-overflow: ellipsis;
    overflow: hidden;
    white-space: nowrap;
    max-width: 10em;
  }

  .navbar-header {
    float: none !important;
  }

  .navbar-toggle {
    display: block !important;
  }

  .tsld.btn-group {
    margin-top: 14px;
  }

  .tsld.btn-group .btn:focus {
    outline: none;
  }

  .title-container {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 940;
    overflow: auto;
    overflow-y: scroll;
  }

  .title {
    position: relative;
    padding-top: 30px;
    padding-bottom: 30px;
    z-index: 950;
    padding: 10px;
    margin-left: auto;
    margin-right: auto;
    text-align: center;
  }

  .subtitle-container {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 940;
    overflow: auto;
    overflow-y: scroll;
  }

  .subtitle {
    position: relative;
    padding-top: 30px;
    padding-bottom: 30px;
    margin-bottom: 5px;
    z-index: 950;
    padding: 10px;
    margin-left: auto;
    margin-right: auto;
    text-align: center;
  }
</style>
