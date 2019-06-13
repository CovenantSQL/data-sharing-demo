"use strict";
/*jslint browser: true, nomen: true*/
/*global $, define, d3, playback*/

define(["./model/model", "./layout/layout", "./frames/replication"], function (Model, Layout, replication) {
  window.archPlayer = playback.player();
  window.replication = replication;
  window.Layout = Layout;
  window.Model = Model;

  // window.archPlayer.frame("replication", "Log Replication", window.replication);
});
