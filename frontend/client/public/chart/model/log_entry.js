
"use strict";
/*jslint browser: true, nomen: true*/
/*global define, d3, tsld*/

define([], function () {
    function LogEntry(model, index, term, command, callback) {
        playback.DataObject.call(this, model);
        this.index = index;
        this.term = term;
        this.command = command;
        this.callback = (callback !== undefined ? callback : null);
    }

    LogEntry.prototype = new playback.DataObject();
    LogEntry.prototype.constructor = LogEntry;

    /**
     * Determines the bounding box of the log entry.
     */
    LogEntry.prototype.bbox = function () {
        return tsld.bbox(this.y, this.x + this.w, this.y + this.h, this.x);
    };

    /**
     * Determines the bounding box of the log entry.
     */
    LogEntry.prototype.bbox = function () {
        return tsld.bbox(this.y, this.x + this.w, this.y + this.h, this.x);
    };

    /**
     * Applies the log to a node.
     */
    LogEntry.prototype.applyTo = function (node) {
        if (node._value <= 0) {
            node._value = 0
        }
        var m = this.command.match(/^(\w+) .*$/);
        switch (m[1].toUpperCase()) {
        case "UPDATE":
            node._value += parseInt(1, 10);
            break;
        case "DELETE":
            node._value += parseInt(1, 10);
            break;
        case "INSERT":
            node._value += parseInt(1, 10);
            break;
        }

        if (this.callback !== null) {
            this.callback();
        }
    };

    LogEntry.prototype.clone = function (model) {
        var clone = new LogEntry(model, this.index, this.term, this.command);
        clone.index = this.index;
        clone.term = this.term;
        clone.command = this.command;
        return clone;
    };

    return LogEntry;
});
