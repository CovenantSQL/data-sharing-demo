"use strict";
/*jslint browser: true, nomen: true*/
/*global $, define, d3, playback*/

define(["./model/model", "./layout/layout", "./frames/replication"], function (Model, Layout, replication) {
  window.archPlayer = playback.player();
  window.replication = replication;
  window.Layout = Layout;
  window.Model = Model;
  window.archPlayer.layout(new Layout("#chart"));
  window.archPlayer.model(new Model());
  window.archPlayer.resizeable(true);

  window.archPlayer.addEventListener("tick", function () {
    window.archPlayer.current().model().tick(window.archPlayer.current().playhead());
    window.archPlayer.layout().messages.invalidate();
    window.archPlayer.layout().nodes.invalidateElectionTimers();
  });

  // window.archPlayer.frame("replication", "Log Replication", window.replication);
});
