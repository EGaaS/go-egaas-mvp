var CodeGenerator = {};

CodeGenerator.Controller = JS_CLASS({
    constructor: function (param) {
        CP(this, param);

        this.model = new CodeGenerator.Model();
        this.over = new CodeGenerator.Over({
            $over: this.$containerWrapper.find(".js-over"),
            model: this.model,
            owner: this
        });

        //this.$over = this.$containerWrapper.find(".js-over");
        //this.$overInner = this.$over.find(".b-over__inner");
        //$(document.body).append(this.$over);
        this.over.hide();

        this.$bag = this.$containerWrapper.find(".js-draggable-bag");
        this.$bagInner = this.$bag.find(".js-draggable-bag__inner");

        this.$trash = $(".js-trash");
        this.$undo = $(".js-undo");
        this.$redo = $(".js-redo");

        this.instrumentPanel = new InstrumentPanel({
            $instrumentPanel: $(".js-instrument-panel")
        });


        this.events();

    },

    setJsonData: function (json) {
        this.model.setJsonData(json);
    },

    generateCode: function () {
        this.code = (new MainTemplate(this.model.json)).renderCode();
        this.$codeGenerated.html(this.code);
    },

    getCode: function () {
        return this.code;
    },

    generateHTML: function () {
        return (new MainTemplate(this.model.json)).renderHTML();
    },

    printJSON: function () {
        this.$output.html(JSON.stringify(this.model.json));
    },

    createDraggableBag: function (html) {
        //this.$bag = $('<div class="b-draggable-bag"><div class="b-draggable-bag__inner">'+html+'</div></div>');
        //this.$bag.hide();
        //$(document.body).append(this.$bag);
        this.$bagInner.html(html);

    },

    removeDraggableBag: function () {
        if(this.$bag)
            this.$bag.hide();
    },

    moveDraggableBag: function (x, y) {
        this.$bag
            .show()
            .css("top", y + 3 - $(window).scrollTop())
            .css("left", x + 3 - $(window).scrollLeft());
    },

    events: function () {
        var self = this;

        this.recalcContainerOffset();

        this.$containerWrapper
            .on("mousemove", function (e) {

                if(self.findOverTagTimer) {
                    return;
                }

                var gap = 3;
                if(self.startDraggingX && self.startDraggingY && !(e.pageX > self.startDraggingX - gap && e.pageX < self.startDraggingX + gap
                    && e.pageY > self.startDraggingY - gap && e.pageY < self.startDraggingY + gap)) {
                    self.startDragging();
                }

                self.findOverTagTimer = setTimeout((function() {
                    self.findOverTagTimer = null;
                }).bind(self), 50);

                self.over.findOverTag(
                    e.pageX - self.containerOffset.left,
                    e.pageY - self.containerOffset.top
                );



                console.log("overTag", self.over.tag);
                console.log("parentOverTag", self.over.parentTag);

                if(self.dragging) {
                    self.over.mode = "drag";
                }
                else {
                    self.over.mode = "view";
                }
                self.over.draw();

            })
            .on("mouseleave", function () {
                self.over.hide();
            });

        this.over.$over
            .on("mousedown", function (e) {
                if(self.over.tag && self.over.tag.type == "tag") {
                    self.draggingTag = self.over.tag;
                    self.$draggingTag = self.$container.find("*[tag-id=" + self.over.tag.id + "]");
                    //console.log("self.$draggingTag", self.$draggingTag);
                    if(self.$draggingTag && self.$draggingTag.length) {
                        self.preStartDragging(e);
                    }
                }
            })
            .on("mouseup", function () {
                if(self.dragging) {
                    self.cancelDragging();
                    self.dropTo(self.over.tag);
                }
                else {
                    self.preStopDragging();
                }
            });

        $(document.body)
            .off("mouseup", function(e) {
                self.onBodyMouseUp(e);
            })
            .on("mouseup", function (e) {
                self.onBodyMouseUp(e);
            })
            .off("mouseleave", function(e) {
                self.onBodyMouseLeave(e);
            })
            .on("mouseleave", function (e) {
                self.onBodyMouseLeave(e);
            })
            .off("mousemove", function(e) {
                self.onBodyMouseMove(e)
            })
            .on("mousemove",  function(e) {
                self.onBodyMouseMove(e)
            });




        // $(".b-over").on("mouseout", function () {
        //     $(this).hide();
        // });

        $(".js-draggable")
            .off("dragstart")
            .on("dragstart", function (e) {
                e.preventDefault();
            })
            .off("mousedown")
            .on("mousedown", function (e) {
                e.preventDefault();
                e.stopPropagation();
                self.$draggingTag = $(this);
                self.draggingTag = {
                    "type": "tag",
                    "name": $(this).attr("tag-name"),
                    "body": []
                };

                if($(this).attr("tag-params")) {
                    self.draggingTag.params = JSON.parse($(this).attr("tag-params"));
                }

                if($(this).attr("tag-body")) {
                    self.draggingTag.body = JSON.parse($(this).attr("tag-body"));
                }
                // if($(this).attr("tag-nestedClassList")) {
                //     self.draggingTag.nestedClassList = JSON.parse($(this).attr("tag-nestedClassList"));
                // }
                //ol, liClass

                self.startDragging(e);
            })
            .on("click", function (e) {
                e.preventDefault();
            });

        this.$trash
            .on("mouseover", function (e) {
                if(self.dragging && self.draggingTag && self.draggingTag.id) {
                    $(this).addClass("b-trash_over");
                }
            })
            .on("mouseout", function (e) {
                $(this).removeClass("b-trash_over");
            })
            .on("mouseup", function (e) {
                if(self.dragging && self.draggingTag && self.draggingTag.id) {
                    self.cancelDragging();
                    self.model.remove(self.draggingTag);
                    self.generateCode();
                    self.render();
                }
            });

        this.$undo.on("click", function () {
            if(self.model.undo()) {
                self.generateCode();
                self.render();
            }
        });

        this.$redo.on("click", function () {
            if(self.model.redo()) {
                self.generateCode();
                self.render();
            }
        });

    },

    onBodyMouseMove: function (e) {
        if(!this.dragging)
            return;
        this.moveDraggableBag(e.pageX, e.pageY);
    },

    onBodyMouseUp: function (e) {
        if(!this.dragging)
            return;
        this.cancelDragging();
    },

    onBodyMouseLeave: function () {
        this.cancelDragging();
    },

    preStartDragging: function (e) {
        this.startDraggingX = e.pageX;
        this.startDraggingY = e.pageY;
        console.log("preStartDragging", e);
    },

    preStopDragging: function () {
        this.startDraggingX = null;
        this.startDraggingY = null;
    },

    startDragging: function () {
        this.dragging = true;
        //this.$overInner.addClass("g-move");
        $(document.body)
            .addClass("g-no-select")
            .addClass("g-no-overflow-x");
        var html = this.$draggingTag.get(0).outerHTML;

        if(this.$draggingTag.find(".js-source-element__preview").length > 0)
            html = this.$draggingTag.find(".js-source-element__preview").get(0).innerHTML;

        this.createDraggableBag(html);
        this.$draggingTag.addClass("b-draggable_dragging");
    },

    cancelDragging: function () {
        this.dragging = false;
        this.preStopDragging();
        //this.$overInner.removeClass("g-move");
        this.over.$over
            .removeClass("b-droppable_inside")
            .removeClass("b-droppable_before")
            .removeClass("b-droppable_after");
        $(document.body)
            .removeClass("g-no-select")
            .removeClass("g-no-overflow-x");
        this.removeDraggableBag();
        this.startDraggingX = null;
        this.startDraggingY = null;


        $(".b-draggable_dragging").removeClass("b-draggable_dragging");
    },

    dropTo: function (overTag) {
        //self.$draggingTag = $("*[tag-id=" + self.over.tag.id + "]");

        console.log("dropTo", overTag, this.over.position, "what", this.draggingTag);
        if(!this.over.canDrop)
            return;
        //if(this.over.position == "inside") {
        //    this.json
        //}
        if(overTag && this.draggingTag && overTag.id && overTag.id === this.draggingTag.id)
            return;
        //console.log("drop", overTag, this.draggingTag, overTag.id, this.draggingTag.id, overTag.id === this.draggingTag.id);
        this.model.appendToTree(this.draggingTag, overTag, this.over.position);

        this.generateCode();
        this.render();
    },

    setHTML: function (html) {
        this.html = html;
    },

    render: function () {
        var self = this;
        this.$container.html(this.generateHTML());
        //alert("rendered");
        this.recalcTagsCoords();
    },

    recalcContainerOffset: function () {
        this.containerOffset = this.$container.offset();
    },

    recalcTagsCoords: function () {
        var self = this;
        var $imgs = this.$container.find("img");
        var loadCounter = 0;
        if($imgs.length == 0)
            this.calcTagsCoords();
        this.$container.find("img").each(function() {
            $(this).on("load error", function () {
                loadCounter++;
                if(loadCounter >= $imgs.length)
                    self.calcTagsCoords();
            });
        });
    },

    calcTagsCoords: function () {
        this.recalcContainerOffset();
        console.log("cont offset", this.containerOffset.top, this.containerOffset.left);
        this.calcTagCoords(this.model.json);
    },

    calcTagCoords: function (tag) {
        if(tag.type == "Template") {
            tag.coords = {
                "left": 0,
                "top": 0,
                "width": "100%",
                "height": "100%"
            };
        }
        if(tag.id) {
            var $tag = this.$container.find('*[tag-id="'+tag.id+'"]');
            var offset = $tag.offset();
            //console.log(tag.name + " offset", '*[tag-id="'+tag.id+'"]', $tag, offset);
            if(offset) {
                tag.coords = {
                    left: offset.left - this.containerOffset.left,
                    top: offset.top - this.containerOffset.top,
                    width: $tag.outerWidth(),
                    height: $tag.outerHeight()
                };
            }
        }
        if(tag.body) {
            for (var i = 0; i < tag.body.length; i++) {
                this.calcTagCoords(tag.body[i]);
            }
        }
    },


});

CodeGenerator.Model = JS_CLASS({
    constructor: function (param) {
        this.history = [];
        this.currentHistoryPos = -1;
        this.json = null;
        CP(this, param);
    },

    saveHistory: function () {
        if(this.currentHistoryPos < this.history.length - 1) {

            this.history.splice(this.currentHistoryPos + 1, 10000);
            console.log("history this.history.splice(" + (this.currentHistoryPos + 1) + ", 10000)");
        }
        this.history.push($.extend(true, {}, this.json));
        this.currentHistoryPos = this.history.length - 1;
        console.log("history", this.history, this.currentHistoryPos);
    },

    undo: function () {
        if(this.currentHistoryPos > 0) {
            this.currentHistoryPos--;
            this.json = $.extend(true, {}, this.history[this.currentHistoryPos]);
            console.log("history undo", this.history, this.currentHistoryPos);
            return true;
        }
        return false;
    },

    redo: function () {
        if(this.currentHistoryPos < this.history.length - 1) {
            this.currentHistoryPos++;
            this.json = $.extend(true, {}, this.history[this.currentHistoryPos]);
            console.log("history redo", this.history, this.currentHistoryPos);
            return true;
        }
        return false;
    },

    setJsonData: function (json) {
        this.json = json;
        this.saveHistory();
    },

    setIds: function (tag) {
        if(!tag)
            return;
        if(!tag.id) {
            tag.id = this.generateId();
        }

        if(tag.body) {
            for (var i = 0; i < tag.body.length; i++) {
                this.setIds(tag.body[i]);
            }
        }
    },

    appendToTree: function (tag, toTag, position) {
        var move = false;
        var inserted = false;
        if(toTag) {

            if(tag.id) {
                console.log("move");
                move = true;
            }
            else {
                console.log("new");
                move = false;
                //tag.id = this.generateId();
                this.setIds(tag);

            }

            if(move) {
                //удаляем из предыдущего пложения
                this.findElementById(this.json, tag.id);
                console.log("remove from prev pos", this.findInfo);

                if (this.findInfo.el && this.findInfo.parent) {
                    this.findInfo.parent.body.splice(this.findInfo.parentPosition, 1);
                    console.log("this.findInfo.parent.body.splice(" + this.findInfo.parentPosition + ", 1)");
                    console.log("parent: ", this.findInfo.parent);
                }
            }


            this.findElementById(this.json, toTag.id);
            if (this.findInfo.el) {
                console.log("appendToTree found", this.findInfo);
                if(position == "inside") {
                    this.findInfo.el.body.push(tag);
                    //inserted = true;
                }
                if(position == "before") {
                    var newPosition = this.findInfo.parentPosition - 1;
                    if(newPosition < 0)
                        newPosition = 0;
                    this.findInfo.parent.body.splice(newPosition, 0, tag);
                    //inserted = true;
                }
                if(position == "after") {
                    var newPosition = this.findInfo.parentPosition + 1;
                    this.findInfo.parent.body.splice(newPosition, 0, tag);
                    //inserted = true;
                }
            }

            this.saveHistory();


        }
    },

    remove: function (tag) {
        this.findElementById(this.json, tag.id);
        console.log("remove tag", this.findInfo);
        if (this.findInfo.el && this.findInfo.parent) {
            this.findInfo.parent.body.splice(this.findInfo.parentPosition, 1);
            this.saveHistory();
        }
    },

    // update: function (id, params) {
    //     var tag = this.findElementById(this.json, id);
    // },

    findElementById: function (el, id) {
        this.findInfo = {
            el: null,
            parent: null,
            parentPosition: 0
        };
        console.log("findElementById", el, id);
        this.findNextElementById(el, id);
    },

    findNextElementById: function (el, id) {


        if (el.id == id) {
            this.findInfo.el = el;
            return;
        }
        if (el.body) {

            for (var i = 0; i < el.body.length; i++) {
                if(this.findInfo.el)
                    break;
                this.findInfo.parent = el;
                this.findInfo.parentPosition = i;
                this.findNextElementById(el.body[i], id);
            }
        }
    },

    generateId: function() {
        return "tag_" + (10000000 + Math.floor(Math.random() * 89999999));
    }
});

CodeGenerator.Over = JS_CLASS({
    constructor: function (param) {
        CP(this, param);
        this.$info = this.$over.find(".js-over-info");
        this.$settings = this.$over.find(".js-over-settings");
        this.$remove = this.$over.find(".js-over-remove");
        this.tag = null;
        this.tagObj = null;
        this.prevTag = null;
        this.parentTag = null;
        this.position = null; //inside, before, after
        this.canDrop = false;
        this.mode = "view"; //"drag"
        this.events();
    },

    events: function () {
        var self = this;
        this.$over.on("click", function() {
            if(self.owner.dragging)
                return;
            if(!self.tag.name)
                return;
            console.log("settings", self.tag);
            new TagSettingsPanel({
                tag: self.tag,
                model: self.model,
                owner: self
            });
        });

        this.$remove.on("click", function(e) {
            e.stopPropagation();
            if(self.owner.dragging)
                return;
            if(!self.tag || !self.tag.name)
                return;

            self.model.remove(self.tag);
            self.owner.generateCode();
            self.owner.render();

            self.owner.$containerWrapper.trigger("mousemove", e);

        });
    },

    hide: function () {
        this.$over.hide();
    },

    findOverTag: function (x, y) {
        this.tag = null;
        this.parentTag = null;
        this.tmpParentOverTag = null;

        this.position = "inside";
        this.findNextOverTag(x, y, this.model.json);

        if(!this.tag) {
            this.tag = this.model.json;
            this.$settings.hide();
            this.$remove.hide();
            // this.tag.coords = {
            //     "left": 0,
            //         "top": 0,
            //         "width": "100%",
            //         "height": "100%"
            // };
        }
        else {
            this.$settings.show();
            this.$remove.show();
        }
        if(this.prevTag != this.tag) {
            this.onOverTagChange();
        }
        this.prevTag = this.tag;

        if(this.tag) {
            this.tagObj = constructTag(this.tag);
        }
        else {
            this.tagObj = null;
        }
        //return this.tag;
    },
    findNextOverTag: function (x, y, tag) {
        if(tag) {
            if (tag.id && tag.coords) {
                if (x >= tag.coords.left && x <= tag.coords.left + tag.coords.width
                    && y >= tag.coords.top && y <= tag.coords.top + tag.coords.height) {

                    if (!this.tag || tag.coords.width < this.tag.coords.width || tag.coords.height < this.tag.coords.height) {
                        this.tag = tag;
                        this.parentTag = this.tmpParentOverTag;

                        this.position = "inside";

                        if (x < tag.coords.left + tag.coords.width * 0.33
                            && y < tag.coords.top + tag.coords.height * 0.33) {
                            this.position = "before";
                        }

                        if (x > tag.coords.left + tag.coords.width * (1 - 0.33)
                            && y > tag.coords.top + tag.coords.height * (1 - 0.33)) {
                            this.position = "after";
                        }
                    }
                }
            }
            if (tag.body) {
                for (var i = 0; i < tag.body.length; i++) {
                    this.tmpParentOverTag = tag;
                    this.findNextOverTag(x, y, tag.body[i]);
                }
            }
        }
    },

    onOverTagChange: function () {

    },

    draw: function () {
        if(this.mode == "drag")
            this.drawDrop();

        if(this.mode == "view")
            this.drawView();

        if(this.tagObj) {
            var name = "Container";
            if(this.tagObj.name) {
                name = this.tagObj.title + ' ' + ((this.tagObj.name != this.tagObj.title) ? this.tagObj.name : "");
            }

            this.$info
                .show()
                .html("<b>" + name + "</b>");
        }
        else
            this.$info.hide();

    },

    drawDrop: function () {
        if(this.tag) {
            this.canDrop = false;
            //если пытаемся положить внутрь перетаскиваемого тега, не пускаем
            if (this.owner.$container.find("*[tag-id=" + this.tag.id + "]").closest(".b-draggable_dragging").length) {
                return;
            }

            this.$over
                .removeClass("b-droppable_inside")
                .removeClass("b-droppable_before")
                .removeClass("b-droppable_after");

            if (this.position == "inside" && this.tag.type != "Template") {
                var tag = constructTag(this.tag);
                if (!tag.accept(this.owner.draggingTag.name))
                    this.position = "after";
                //console.log("tag", tag);
            }

            //в основной контейнер - только внутрь
            if (this.tag.type == "Template") {
                this.position = "inside";
                var tag = constructTag(this.tag);
                if (!tag.accept(this.owner.draggingTag.name))
                    return;
            }

            if (this.position == "before" || this.position == "after") {
                //before/after - проверка accept parent
                if (this.parentTag) {
                    var parentTag = constructTag(this.parentTag);
                    if (!parentTag   //в самый высокий уровень можно только класть внутрь, а не до/после
                        || !parentTag.accept(this.owner.draggingTag.name)) {
                        return;
                    }
                }
            }

            if(this.tag.coords) {
                this.$over
                    .show()
                    .css("left", this.tag.coords.left)
                    .css("top", this.tag.coords.top)
                    .css("width", this.tag.coords.width)
                    .css("height", this.tag.coords.height)
                    .addClass("b-droppable_" + this.position);
            }

            this.canDrop = true;
        }
    },

    drawView: function () {
        if(this.tag) {
            if(this.tag.coords) {
                this.$over
                    .show()
                    .css("left", this.tag.coords.left)
                    .css("top", this.tag.coords.top)
                    .css("width", this.tag.coords.width)
                    .css("height", this.tag.coords.height);
            }
        }
    }

});

var TagSettingsPanel = JS_CLASS({

    constructor: function (param) {
        CP(this, param);

        this.$panel = $(".js-tag-settings-panel");
        this.$panelBody = this.$panel.find(".js-body");
        this.$controls = this.$panel.find(".js-controls");
        this.$close = this.$panel.find(".js-close");
        this.$save = this.$panel.find(".js-save");
        this.$title = this.$panel.find(".js-title");

        this.tagObj = constructTag(this.tag);

        this.controls = [];

        this.events();
        this.show();
        this.render();
    },

    events: function () {
        var self = this;
        // this.$panel.on("click", function () {
        //     self.cancel();
        // });

        this.$panelBody.on("click", function (e) {
            e.stopPropagation();
        });

        this.$close.on("click", function () {
            self.cancel();
        });

        this.$save.on("click", function () {
            self.save();
        });
    },

    show: function () {
        this.$panel.show();
    },

    cancel: function () {
        this.$panel.hide();
    },

    save: function () {

        for(var i = 0; i < this.controls.length; i++) {
            var control = this.controls[i];
            console.log(control.name, control.getValue());
            this.tag.params[control.name] = control.getValue();
        }

        //this.$panel.hide();
        this.owner.owner.generateCode();
        this.owner.owner.render();
        this.model.saveHistory();
    },

    render: function () {
        this.$controls.empty();
        if(!this.tagObj)
            return;

        this.$title.html(this.tagObj.title);

        for(var paramName in this.tagObj.paramsType) {
            if (this.tagObj.paramsType.hasOwnProperty(paramName)) {
                //paramArr.push(paramName + ' = "' + this.params[paramName] + '"');
                var param = this.tagObj.paramsType[paramName];
                if(Control[param.type]) {
                    var control = new Control[param.type]({
                        $content: this.$controls,
                        name: paramName,
                        param: param
                    });
                    control.setValue(this.tagObj.params[paramName]);
                    this.controls.push(control);
                }
            }
        }

    }


});


function constructTag(item, offset) {
    var tagObj = null;
    if(item.type == "Template") {
        tagObj = new MainTemplate(item, offset);
    }

    if(window["Tag" + item.name])
        tagObj = new window["Tag" + item.name] (item, offset);

    // switch(item.name) {
    //     case "A":
    //         tagObj = new TagA(item, offset);
    //         break;
    //     case "P":
    //         tagObj = new TagP(item, offset);
    //         break;
    //     case "Div":
    //         tagObj = new TagDiv(item, offset);
    //         break;
    //     case "Image":
    //         tagObj = new TagImage(item, offset);
    //         break;
    //     case "Li":
    //         tagObj = new TagLi(item, offset);
    //         break;
    //     case "Divs":
    //         tagObj = new TagDivs(item, offset);
    //         break;
    //     case "UList":
    //         tagObj = new TagUList(item, offset);
    //         break;
    //     case "LiBegin":
    //         tagObj = new TagLiBegin(item, offset);
    //         break;
    // }
    return tagObj;
}

var Tag = JS_CLASS({
    acceptRule: null,
    exceptRule: null,
    paramsType: {},
    params: {},
    constructor: function (param) {
        CP(this, param);
    },
    renderOffset: function () {
        if(this.lineOffset)
            return Array((this.lineOffset) * 2).join(" ");
        return "";
    },
    renderHTML: function () {
        return "";
    },
    accept: function (tagName) {

        if(this.exceptRule) {
            var exceptRuleArr = this.exceptRule.split(" ");
            if ($.inArray(tagName, exceptRuleArr) > -1)
                return false;
        }

        if(!this.acceptRule)
            return false;
        if(this.acceptRule == "*")
            return true;
        var acceptRuleArr = this.acceptRule.split(" ");
        return ($.inArray(tagName, acceptRuleArr) > -1);
    }
});

/*простые теги, например a, p, img*/
var SimpleTag = JS_CLASS(Tag, {
    name: "",
    lineOffset: 0,
    acceptRule: null,
    constructor: function (param, lineOffset) {
        CP(this, param);
        if(lineOffset)
            this.lineOffset = lineOffset;
        console.log("construct " + this.name, this.lineOffset);
    },
    renderCode: function (named) {

        var code = this.renderOffset() + this.name + (named ? "{" : "(");
        var paramArr = [];
        for(var paramName in this.paramsType) {
            if (this.paramsType.hasOwnProperty(paramName)) {
                var value = "";
                if(typeof this.params[paramName] !== "undefined" && this.params[paramName] !== null)
                    value = this.params[paramName];
                if(named)
                    paramArr.push(paramName + ' = "' + value + '"');
                else {
                    paramArr.push('"' + value + '"');
                }
            }
        }
        code += paramArr.join(", ");
        code += (named ? "}" : ")") + "\n";
        return code;
    }
});

/* теги со вложенностью, например LiBegin, UList*/
var StructureTag = JS_CLASS(Tag, {
    nameBegin: "",
    nameEnd: "",
    lineOffset: 0,
    acceptRule: "*",
    exceptRule: "Li LiBegin",
    constructor: function (param, lineOffset) {
        CP(this, param);
        if(lineOffset)
            this.lineOffset = lineOffset;
    },
    renderCode: function (named) {
        var code = this.renderOffset() + this.nameBegin + (named ? "{" : "(");
        var paramArr = [];
        for(var paramName in this.paramsType) {
            if (this.paramsType.hasOwnProperty(paramName)) {
                var value = "";
                if(typeof this.params[paramName] !== "undefined" && this.params[paramName] !== null)
                    value = this.params[paramName];
                if(named)
                    paramArr.push(paramName + ' = "' + value + '"');
                else {
                    paramArr.push('"' + value + '"');
                }
            }
        }
        code += paramArr.join(", ");
        code += (named ? "}" : ")") + "\n";
        code += this.renderSubItems();

        code += this.renderOffset() + this.nameEnd + ":\n";
        return code;
    },

    renderSubItems: function (returnType) {
        var code = "";
        if(this.body) {
            for(var i = 0; i < this.body.length; i++) {
                var item = this.body[i];

                var tagObj = constructTag(item, this.lineOffset + 1);

                if(tagObj) {
                    if(returnType == 'html') {
                        code += tagObj.renderHTML();
                    }
                    else {
                        code += tagObj.renderCode();
                    }
                }
            }
        }
        return code;
    }

});


var MainTemplate = JS_CLASS(StructureTag, {
    nameBegin: "",
    nameEnd: "",
    renderCode: function () {
        return this.renderSubItems();
    },

    renderHTML: function () {
        return this.renderSubItems('html');
    }

});

var TagA = JS_CLASS(SimpleTag, {
    name: "A",
    title: "Link",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Link text",
            "type": "String",
            "obligatory": true
        },
        "href": {
            "title": "URL",
            "description": "Enter link URL",
            "type": "String",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<a tag-id="' + this.id + '" href="' + this.params.href + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</a>';
        return html;
    }
});

var TagDiv = JS_CLASS(SimpleTag, {
    name: "Div",
    title: "Block",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "String",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<div tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</div>';
        return html;
    }
});

var TagDivs = JS_CLASS(StructureTag, {
    nameBegin: "Divs",
    nameEnd: "DivsEnd",
    title: "Blocks",
    paramsType: {
        "nestedClassList": {
            "title": "Class list",
            "description": "Input CSS class list using whitespace separator, separate class groups with commas. Example: \"row, col-xs-10 col-md-7\"",
            "type": "WhiteSpaceStringArray",
            "obligatory": false
        }
    },
    renderCode: function () {
        var code = this.renderOffset() + this.nameBegin + "(";

        if(this.params.nestedClassList)
            code += this.params.nestedClassList.join(", ");
        code += ")\n";

        code += this.renderSubItems();

        code += this.renderOffset() + this.nameEnd + ":\n";
        return code;
    },

    renderHTML: function () {
        var html = '';

        if(this.params.nestedClassList) {
            for (var i = 0; i < this.params.nestedClassList.length; i++) {
                html += '<div tag-id="' + this.id + '" class="' + this.params.nestedClassList[i] + '">';
            }
        }
        else {
            html += '<div tag-id="' + this.id + '">';
        }

        html += this.renderSubItems('html');

        if(this.params.nestedClassList) {
            for (i = 0; i < this.params.nestedClassList.length; i++) {
                html += '</div>';
            }
        }
        else {
            html += '</div>';
        }

        return html;
    }
});

var TagP = JS_CLASS(SimpleTag, {
    name: "P",
    title: "Text",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<p tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</p>';
        return html;
    }
});

var TagEm = JS_CLASS(SimpleTag, {
    name: "Em",
    title: "Cursive",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<em tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</em>';
        return html;
    }
});

var TagLi = JS_CLASS(SimpleTag, {
    name: "Li",
    title: "List element",
    paramsType: {
        "text": {
            "title": "Text",
            "type": "String",
            "obligatory": true
        },
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<li tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</li>';
        return html;
    }
});

var TagLiBegin = JS_CLASS(StructureTag, {
    nameBegin: "LiBegin",
    nameEnd: "LiEnd",
    title: "Complex list element",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<li tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">';

        html += this.renderSubItems('html');

        html += '</li>';
        return html;
    }
});

var TagSmall = JS_CLASS(SimpleTag, {
    name: "Small",
    title: "Small",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<small tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</small>';
        return html;
    }
});

var TagSpan = JS_CLASS(SimpleTag, {
    name: "Span",
    title: "Span",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<span tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</span>';
        return html;
    }
});

var TagStrong = JS_CLASS(SimpleTag, {
    name: "Strong",
    title: "Bold",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<strong tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</strong>';
        return html;
    }
});

var TagLabel = JS_CLASS(SimpleTag, {
    name: "Label",
    title: "Label",
    paramsType: {
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        },
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<label tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</label>';
        return html;
    }
});

var TagLegend = JS_CLASS(SimpleTag, {
    name: "Legend",
    title: "Legend",
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": true
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<legend tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</legend>';
        return html;
    }
});

var TagTag = JS_CLASS(SimpleTag, {
    name: "Tag",
    title: "Tag",
    paramsType: {
        "tagname": {
            "title": "Select tag name",
            "type": "Select",
            "values": [
                {"name": "Header 1", "value": "h1"},
                {"name": "Header 2", "value": "h2"},
                {"name": "Header 3", "value": "h3"},
                {"name": "Header 4", "value": "h4"},
                {"name": "Header 5", "value": "h5"},
                {"name": "Header 6", "value": "h6"},
                {"name": "Button", "value": "button"}
            ],
            "obligatory": true
        },
        "text": {
            "title": "Text",
            "type": "Textarea",
            "obligatory": false
        },
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<' + this.params.tagname + ' tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + (this.params.text ? this.params.text : "") + '</' + this.params.tagname + '>';
        return html;
    }
});

var TagImage = JS_CLASS(SimpleTag, {
    name: "Image",
    title: "Image",
    paramsType: {
        "src": {
            "title": "Image URL",
            "description": "Example: http://site.com/image.jpg or img/image.jpg",
            "type": "String",
            "obligatory": true
        },
        "alt": {
            "title": "Alternative text",
            "description": "Enter alternative image text",
            "type": "String",
            "obligatory": false
        },
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<img tag-id="' + this.id + '" src="' + (this.params.src ? this.params.src : "") + '" class="' + (this.params.class ? this.params.class : "") + '" alt="' + (this.params.alt ? this.params.alt : "") + '">';
        return html;
    }
});

var TagImageInput = JS_CLASS(SimpleTag, {
    name: "ImageInput",
    title: "ImageInput",
    paramsType: {
        "id": {
            "title": "Image ID",
            "type": "String",
            "obligatory": true
        },
        "width": {
            "title": "Image width",
            "type": "String",
            "obligatory": false
        },
        "ratio_height": {
            "title": "Image height or aspect ratio",
            "description": "Enter image height or aspect ratio, for example 1/2",
            "type": "String",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var html = '<button tag-id="' + this.id + '" class="btn btn-primary">Select image</button>';
        return html;
    }
});

var TagMarkDown = JS_CLASS(SimpleTag, {
    name: "MarkDown",
    title: "MarkDown",
    paramsType: {
        "text": {
            "title": "Text",
            "description": "Text to be converted to HTML",
            "type": "Textarea",
            "obligatory": true
        }
    },
    renderHTML: function () {
        var html = '<span tag-id="' + this.id + '">' + (this.params.text ? this.params.text : "") + '</span>';
        return html;
    }
});


var TagUList = JS_CLASS(StructureTag, {
    nameBegin: "UList",
    nameEnd: "UListEnd",
    title: "List",
    acceptRule: "Li LiBegin",
    exceptRule: null,
    paramsType: {
        "class": {
            "title": "Element class list",
            "description": "Input CSS class list using whitespace separator",
            "type": "WhiteSpaceString",
            "obligatory": false
        },
        "ol": {
            "title": "Numbered list",
            "description": "Check this box for numbered list",
            "type": "Checkbox",
            "onValue": "ol",
            "offValue": "",
            "obligatory": false
        },
        "liClass": {
            "title": "Common class for list elements",
            "description": "Input CSS class list using whitespace separator",
            "type": "String",
            "obligatory": false
        }
    },
    renderHTML: function () {
        var tag = 'ul';
        if(this.params.ol == 'ol')
            tag = 'ol';
        var html = '<' + tag + ' tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">';

        html += this.renderSubItems('html');

        html += '</' + tag + '>';

        return html;
    }
});



var Control = {};

Control.BaseController = JS_CLASS({
    $content: null,
    value: null,
    name: "",
    constructor: function (param) {
        CP(this, param);

        this.data = this.param;
        this.data.name = this.name;

        //console.log("checkbox", this.data);

        this.$html = $(TPL ($(this.tpl).html(), this.data));
        this.$content.append(this.$html);

        this.$input = this.$html.find(".js-control-input");
        this.$error = this.$html.find(".js-control-error");
        this.$description = this.$html.find(".js-control-description");
        if(this.data && this.data.description)
            this.$description.show();

        this.init();
        //console.log(this.data);
    },

    init: function() {

    },

    setValue: function(value) {
        console.log(this.name + " set " + value);
        if(value === null)
            value = '';
        this.$input.val(value);
    },

    getValue: function() {
        return this.$input.val();
    }
});

Control.String = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-string",
    constructor: function (param) {
        SUPER(this,arguments);
    },

    init: function () {

    }
});

Control.Textarea = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-textarea",
    constructor: function (param) {
        SUPER(this,arguments);
    },

    init: function () {

    }
});

Control.WhiteSpaceString = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-white-space-string",
    constructor: function (param) {
        SUPER(this,arguments);
    },

    init: function () {

    }
});

Control.WhiteSpaceStringArray = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-white-space-string-array",
    constructor: function (param) {
        SUPER(this,arguments);
    },

    init: function () {

    },

    setValue: function(value) {

        if(value === null)
            value = '';
        if(value && value.length)
            value = value.join(", ");
        this.$input.val(value);
    },

    getValue: function() {
        var value = this.$input.val();
        if(value.length)
            return value.split(", ");
        return [];
    }
});

Control.Checkbox = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-checkbox",
    constructor: function (param) {

        SUPER(this,arguments);

        if(typeof this.param.onValue == "undefined")
            this.param.onValue = 1;
        if(typeof this.param.offValue == "undefined")
            this.param.offValue = 0;
    },

    init: function () {

    },

    setValue: function(value) {
        if(value === null)
            value = '';
        console.log("checkbox ol", value, this.param.onValue);
        this.$input.prop("checked", value === this.param.onValue);
    },

    getValue: function(value) {
        return this.$input.prop("checked") ? this.param.onValue : this.param.offValue;
    }
});

Control.Select = JS_CLASS(Control.BaseController, {
    tpl: "#tpl-control-select",
    constructor: function (param) {
        SUPER(this,arguments);
    },

    init: function () {
        if(this.param.values && this.param.values.length) {
            for(var i = 0; i < this.param.values.length; i++) {
                this.$input.append('<option value="'+this.param.values[i].value+'">'+this.param.values[i].name+'</option>');
            }
        }
    }
});

var InstrumentPanel = JS_CLASS({
    elements: [
        {
            "title": "Link",
            "type": "tag",
            "name": "A",
            "params": {
                "text": "Link URL",
                "href": "http://"
            }
        },
        {
            "title": "Text Block",
            "type": "tag",
            "name": "Div",
            "params": {
                "text": "Block"

            }
        },
        {
            "title": "Nested Blocks",
            "type": "tag",
            "name": "Divs",
            "params": {
            }
        },
        {
            "title": "Panel",
            "type": "tag",
            "name": "Divs",
            "params": {
                "nestedClassList": [
                    "panel panel-primary"
                ]
            },
            "body": [
                {
                    "type": "tag",
                    "name": "Divs",
                    "params": {
                        "nestedClassList": [
                            "panel-heading"
                        ]
                    },
                    "body": [
                        {
                            "type": "tag",
                            "name": "Div",
                            "params": {
                                "class": "panel-title",
                                "text": "Title"
                            }
                        }
                    ]
                },
                {
                    "type": "tag",
                    "name": "Divs",
                    "params": {
                        "nestedClassList": [
                            "panel-body"
                        ]
                    },
                    "body": [
                        {
                            "type": "tag",
                            "name": "P",
                            "params": {
                                "text": "Panel content"
                            }
                        }
                    ]
                },
                {
                    "type": "tag",
                    "name": "Divs",
                    "params": {
                        "nestedClassList": [
                            "panel-footer"
                        ]
                    },
                    "body": [
                        {
                            "type": "tag",
                            "name": "P",
                            "params": {
                                "text": "Footer"
                            }
                        }
                    ]
                }
            ]

        },

        {
            "title": "Text",
            "type": "tag",
            "name": "P",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Cursive Text",
            "type": "tag",
            "name": "Em",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Small Text",
            "type": "tag",
            "name": "Small",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Span",
            "type": "tag",
            "name": "Span",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Bold Text",
            "type": "tag",
            "name": "Strong",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Label",
            "type": "tag",
            "name": "Label",
            "params": {
                "text": "Label name",
                "class": ""
            }
        },

        {
            "title": "Legend",
            "type": "tag",
            "name": "Legend",
            "params": {
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "H1-H6 Header or Button",
            "type": "tag",
            "name": "Tag",
            "params": {
                "tagname": "h1",
                "class": "",
                "text": "Sample text"
            }
        },

        {
            "title": "Image",
            "type": "tag",
            "name": "Image",
            "params": {
                "src": "http://apla.io/images/prospects_logo.png",
                "alt": "",
                "class": ""
            }
        },

        {
            "title": "Image Upload button",
            "type": "tag",
            "name": "ImageInput",
            "params": {
                "id": "test",
                "width": "100",
                "ratio_height": "1/1"
            }
        },

        {
            "title": "MarkDown HTML ",
            "type": "tag",
            "name": "MarkDown",
            "params": {
                "text": "<b>Sample html</b><hr>"
            }
        },

        {
            "title": "Unordered List",
            "type": "tag",
            "name": "UList",
            "params": {
                "ol": ""
            },
            "body": [
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "First"}
                },
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "Second"}
                },
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "Third"}
                }
            ]
        },

        {
            "title": "Ordered List",
            "type": "tag",
            "name": "UList",
            "params": {
                "ol": "ol"
            },
            "body": [
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "First"}
                },
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "Second"}
                },
                {
                    "type":"tag",
                    "name": "Li",
                    "params": {"text": "Third"}
                }
            ]
        },

        {
            "title": "List Item",
            "type":"tag",
            "name": "Li",
            "params": {"text": "New Item"}
        },

        {
            "title": "Complex List Item",
            "type":"tag",
            "name": "LiBegin",
            "params": {"class": ""},
            "body": [
                {
                    "type": "tag",
                    "name": "Image",
                    "params": {
                        "src": "http://apla.io/images/i19.png"
                    }
                },
                {
                    "type": "tag",
                    "name": "Em",
                    "params": {
                        "class": "",
                        "text": "New item block"
                    }
                }
            ]
        },

    ],
    $instrumentPanel: null,
    $sourceElements: null,
    constructor: function (param) {
        CP(this, param);
        this.init();
    },

    init: function () {
        if(this.$instrumentPanel) {
            this.$sourceElements = this.$instrumentPanel.find(".js-source-elements");

            for(var i = 0; i < this.elements.length; i++) {
                var el = this.elements[i];
                var tag = constructTag(el);

                if(!tag)
                    continue;

                var zIndex = this.elements.length - i;

                var html = "<div class='js-draggable b-source-element' tag-name='"+tag.name+"' tag-params='"+ JSON.stringify(tag.params) + "'";
                if(tag.body)
                    html += " tag-body='"+ JSON.stringify(tag.body) +"'";
                html += " style='z-index: "+zIndex+"'>" + tag.title;
                html += "<div class='b-source-element__preview js-source-element__preview'>" + tag.renderHTML() + "</div></div>";

                this.$sourceElements.append(html);
            }

        }
    }
});
