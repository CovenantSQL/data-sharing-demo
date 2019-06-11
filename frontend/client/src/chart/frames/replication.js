
"use strict";
/*jslint browser: true, nomen: true*/
/*global define*/

define([], function () {
    return function (frame) {
        var player = frame.player(),
            layout = frame.layout(),
            model = function() { return frame.model(); },
            client = function(id) { return frame.model().clients.find(id); },
            node = function(id) { return frame.model().nodes.find(id); },
            cluster = function(value) { model().nodes.toArray().forEach(function(node) { node.cluster(value); }); },
            wait = function() { var self = this; model().controls.show(function() { self.stop(); }); },
            subtitle = function(s, pause) { model().subtitle = s + model().controls.html(); layout.invalidate(); if (pause === undefined) { model().controls.show() }; },
            clear = function() { subtitle('', false); },
            removeAllNodes = function() { model().nodes.toArray().forEach(function(node) { node.state("stopped"); }); model().nodes.removeAll(); };

        //------------------------------
        // Title
        //------------------------------
        frame.after(0, function() {
            model().clear();
        })
        .after(0, function () {
            frame.model().title = '<h2 style="visibility:visible">Log Replication</h1>'
                                + '<br/>' + frame.model().controls.html();
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
        })
        .after(100, function () {
            model().forceImmediateLeader();
        })


        //------------------------------
        // Single Entry Replication
        //------------------------------
        .after(300, function () {
            model().clients.create("X");
            subtitle('<h2>First a client sends a change to the leader.</h2>', false);
        })
        .then(wait).indefinite()
        .then(function () {
            client("X").send(model().leader(), "SET 15");
        })
        .after(model().defaultNetworkLatency, function() {
            subtitle('<h2>The change is appended to the leader\'s log...</h2>');
        })
        .at(model(), "appendEntriesRequestsSent", function () {})
        .after(model().defaultNetworkLatency * 0.25, function(event) {
            subtitle('<h2>...then the change is sent to the followers on the next heartbeat.</h2>');
        })
        .after(1, clear)
        .at(model(), "commitIndexChange", function (event) {
            if(event.target === model().leader()) {
                subtitle('<h2>An entry is committed once a majority of followers acknowledge it...</h2>');
            }
        })
        .after(model().defaultNetworkLatency * 0.25, function(event) {
            subtitle('<h2>...and a response is sent to the client.</h2>');
        })
        .after(1, clear)
        .after(model().defaultNetworkLatency, function(event) {
            subtitle('<h2>Now let\'s send a command to increment the value by "2".</h2>');
            client("X").send(model().leader(), "ADD 2");
        })
        .after(1, clear)
        .at(model(), "recv", function (event) {
            subtitle('<h2>Our system value is now updated to "7".</h2>', false);
        })
        .after(1, wait).indefinite()


        player.play();
    };
});
