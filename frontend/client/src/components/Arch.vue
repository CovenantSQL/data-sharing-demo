<template>
  <div class="container">
    <div id="arch">
    </div>
  </div>
</template>

<script>
  import 'jquery'
  import '../chart/model/model'
  import '../chart/layout/layout'
  import '../chart/frames/init'
  import '../chart/scripts/domReady/domReady-2.0.1'
  import 'd3'
  import '../chart/scripts/playback/playback'
  import '../chart/scripts/tsld/tsld'
  export default {
    name: "Arch",
    created() {
      this.initialize()
    },
    methods: {
      initialize() {
        var i, menu, frame,
          player = playback.player();
        player.layout(new Layout("#chart"));
        player.model(new Model());
        player.resizeable(true);
        frames(player);

        // Handle "continue" button click.
        $(doc).on("click", ".btn.resume", function () {
          player.current().model().controls.resume.click();
        });
        // Refresh the messages on every frame.
        player.addEventListener("tick", function () {
          player.current().model().tick(player.current().playhead());
          player.layout().messages.invalidate();
          player.layout().nodes.invalidateElectionTimers();
        });
        // Write out the frames to the menu.
        menu = $("nav .dropdown-menu");
        menu.empty();
        for (i = 0; i < player.frames().length; i += 1) {
          frame = player.frame(i);
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
