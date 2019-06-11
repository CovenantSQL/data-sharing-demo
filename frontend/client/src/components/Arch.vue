<template>
  <div class="container">
    <div id="arch">
      <div class="container">
        <div id="chart">
        </div>
        <v-btn color="primary" @click="load">load</v-btn>
      </div>
    </div>
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
      load: function () {
        if (this.frame_num === 0) {
          window.archPlayer.frame(this.frame_num, "Log Replication" + this.frame_num, this.initFrame());
        } else {
          // window.archPlayer.current().model().clients.create("Y" + this.frame_num);
          // window.archPlayer.frame(this.frame_num, "Log Replication" + this.frame_num, this.newFrame());
          window.archPlayer.next()
          this.newFrame()(window.archPlayer.current())
        }
        this.frame_num += 1;
      },
      initFrame: function () {
        return function (fr) {
          this.currentFrame = fr;
          var player = fr.player(),
            layout = fr.layout(),
            model = function () {
              return fr.model();
            },
            client = function (id) {
              return fr.model().clients.find(id);
            },
            node = function (id) {
              return fr.model().nodes.find(id);
            },
            cluster = function (value) {
              model().nodes.toArray().forEach(function (node) {
                node.cluster(value);
              });
            };

          fr.after(0, function () {
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
            .after(300, function () {
            })
            .after(500, function () {
              client("X").send(model().leader(), "SET 5");
            })
            .after(model().defaultNetworkLatency, function () {
            })
            .at(model(), "appendEntriesRequestsSent", function () {
            })
            .after(model().defaultNetworkLatency * 0.25, function () {
            })
            .at(model(), "commitIndexChange", function (event) {
            })
            .after(model().defaultNetworkLatency * 0.25, function () {
            })
            .after(model().defaultNetworkLatency, function () {
              client("X").send(model().leader(), "ADD 2");
            })
            .at(model(), "recv", function () {
            })


          player.play();
        }
      },
      newFrame: function () {
        return function (fr) {
          var player = fr.player(),
            model = function () {
              return fr.model();
            },
            client = function (id) {
              return fr.model().clients.find(id);
            }

          fr.after(500, function () {
            client("X").send(model().leader(), "SET 5");
          })
            .after(model().defaultNetworkLatency, function () {
            })
            .at(model(), "appendEntriesRequestsSent", function () {
            })
            .after(model().defaultNetworkLatency * 0.25, function () {
            })
            .at(model(), "commitIndexChange", function (event) {
            })
            .after(model().defaultNetworkLatency * 0.25, function () {
            })
            .after(model().defaultNetworkLatency, function () {
              client("X").send(model().leader(), "ADD 2");
            })
            .at(model(), "recv", function () {
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

  nav {
    z-index: 1000;
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
