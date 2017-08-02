var CodeGenerator = JS_CLASS({

    constructor: function (param) {
        this.json = null;
        CP(this, param);

        this.$over = this.$containerWrapper.find(".js-over");
        //$(document.body).append(this.$over);
        this.$over.hide();

        this.$bag = this.$containerWrapper.find(".js-draggable-bag");
        this.$bagInner = this.$bag.find(".js-draggable-bag__inner");

        this.events();
        //this.eventsDynamic();

    },

    setJsonData: function (json) {
        this.json = json;
    },

    generateCode: function () {
        return (new MainTemplate(this.json)).renderCode();
    },

    generateHTML: function () {
        return (new MainTemplate(this.json)).renderHTML();
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
            .css("top", y + 3)
            .css("left", x + 3);
    },

    events: function () {
        var self = this;

        this.recalcContainerOffset();

        this.$containerWrapper
            .on("mousemove", function (e) {

                if(self.findOverTagTimer) {
                    return;
                }

                self.findOverTagTimer = setTimeout((function() {
                    self.findOverTagTimer = null;
                }).bind(self), 100);

                self.overTag = self.findOverTag(
                    e.pageX - self.containerOffset.left,
                    e.pageY - self.containerOffset.top
                );
                //console.log(e.pageX, e.pageY);
                //console.log("overTag", overTag);
                if(self.overTag) {
                    self.$over
                        .show()
                        .css("left", self.overTag.coords.left)
                        .css("top", self.overTag.coords.top)
                        .css("width", self.overTag.coords.width)
                        .css("height", self.overTag.coords.height)
                        .removeClass("b-over_inside")
                        .removeClass("b-over_before")
                        .removeClass("b-over_after");
                    if(self.dragging) {
                        self.$over.addClass("b-over_" + self.overPosition);
                    }
                }
            })
            .on("mouseleave", function () {
                self.$over.hide();
            });

        this.$over
            .on("mousedown", function () {
                if(self.overTag) {
                    self.$draggingTag = $("*[tag-id=" + self.overTag.id + "]");
                    //console.log("self.$draggingTag", self.$draggingTag);
                    if(self.$draggingTag && self.$draggingTag.length) {
                        self.dragging = true;
                        self.$over.addClass("g-move");
                        $(document.body)
                            .addClass("g-no-select")
                            .addClass("g-no-overflow-x");
                        self.createDraggableBag(self.$draggingTag.get(0).outerHTML);
                        self.$draggingTag.addClass("b-draggable_dragging");
                    }
                }
            })
            .on("mouseup", function () {
                if(!self.dragging)
                    return;
                self.dragging = false;
                self.$over.removeClass("g-move");
                $(document.body)
                    .removeClass("g-no-select")
                    .removeClass("g-no-overflow-x");
                self.removeDraggableBag();
                $(".b-draggable_dragging").removeClass("b-draggable_dragging");
            });

        if(!window.generatorEventsInited) {
            window.generatorEventsInited = true;
            $(document.body)
                .on("mouseup mouseleave", function (e) {
                    if(!self.dragging)
                        return;
                    self.dragging = false;
                    self.$over.removeClass("g-move");
                    $(document.body)
                        .removeClass("g-no-select")
                        .removeClass("g-no-overflow-x");
                    self.removeDraggableBag();
                    // self.$container.find(".js-droppable")
                    //     .removeClass("b-droppable-inside")
                    //     .removeClass("b-droppable-before")
                    //     .removeClass("b-droppable-after");

                    $(".b-draggable_dragging").removeClass("b-draggable_dragging");
                })
                .on("mousemove", function (e) {
                    if(!self.dragging)
                        return;
                    self.moveDraggableBag(e.pageX, e.pageY);
                });
        }



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
                e.stopPropagation();
                self.dragging = true;
                self.$over.addClass("g-move");
                $(document.body)
                    .addClass("g-no-select")
                    .addClass("g-no-overflow-x");
                self.createDraggableBag(this.outerHTML);
                $(this).addClass("b-draggable_dragging");
            });

    },

    eventsDynamic: function () {
        var self = this;

        this.$container.find(".js-droppable")
            .off("mouseover")
            .on("mouseover", function (e) {
                e.stopPropagation();
                if(self.dragging)
                    return;
                //$(this)
                //    .addClass("b-droppable-inside")
                //    .removeClass("b-droppable-before")
                //    .removeClass("b-droppable-after");



                // $(".b-over")
                //     .show()
                //     .css("left", $(this).offset().left)
                //     .css("top", $(this).offset().top)
                //     .css("width", $(this).outerWidth())
                //     .css("height", $(this).outerHeight());
                //console.log($(this).offset().left, $(this).offset().top);

            })
            .off("mouseout")
            .on("mouseout", function (e) {
                e.stopPropagation();
                if(!self.dragging) {
                    //$(this)
                    //    .removeClass("b-droppable-inside");
                    return;
                }

                $(this)
                    .removeClass("b-droppable-inside")
                    .removeClass("b-droppable-before")
                    .removeClass("b-droppable-after");
            })
            .off("mousemove")
            .on("mousemove", function (e) {
                if(!self.dragging)
                    return;

                e.stopPropagation();

                //console.log("move droppable", e);

                self.moveDraggableBag(e.pageX, e.pageY);


                var width = $(this).outerWidth();
                var height = $(this).outerHeight();

                var offsetX = e.pageX - $(this).offset().left;
                var offsetY = e.pageY - $(this).offset().top;

                // inside a scrolling pane:
                //
                // var x = (evt.pageX - $('#element').offset().left) + self.frame.scrollLeft();
                // var y = (evt.pageY - $('#element').offset().top) + self.frame.scrollTop();

                //console.log("el width", width, height, offsetX, offsetY);

                if(offsetX < width * 0.33 && offsetY < height * 0.33) {
                    $(this)
                        .addClass("b-droppable-before")
                        .removeClass("b-droppable-inside")
                        .removeClass("b-droppable-after");
                }
                else {
                    if (offsetX > width * (1 - 0.33) && offsetY > height * (1 - 0.33)) {
                        $(this)
                            .addClass("b-droppable-after")
                            .removeClass("b-droppable-inside")
                            .removeClass("b-droppable-before");
                    }
                    else {
                        $(this)
                            .addClass("b-droppable-inside")
                            .removeClass("b-droppable-before")
                            .removeClass("b-droppable-after");
                    }
                }
            });


    },

    setHTML: function (html) {
        this.html = html;
    },

    render: function () {
        var self = this;
        this.$container.html(this.generateHTML());

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

    recalcContainerOffset: function () {
        this.containerOffset = this.$container.offset();
    },

    calcTagsCoords: function () {
        this.recalcContainerOffset();
        console.log("cont offset", this.containerOffset.top, this.containerOffset.left);
        this.calcTagCoords(this.json);
    },

    calcTagCoords: function (tag) {
        if(tag.id) {
            var $tag = $('[tag-id="'+tag.id+'"]');
            var offset = $tag.offset();
            tag.coords = {
                left: offset.left - this.containerOffset.left,
                top: offset.top - this.containerOffset.top,
                width: $tag.outerWidth(),
                height: $tag.outerHeight()
            };
        }
        if(tag.body) {
            for (var i = 0; i < tag.body.length; i++) {
                this.calcTagCoords(tag.body[i]);
            }
        }
    },

    findOverTag: function (x, y) {
        this.overTag = null;
        this.overPosition = "inside";
        this.findNextTag(x, y, this.json);
        return this.overTag;
    },
    findNextTag: function (x, y, tag) {
        if(tag) {
            if (tag.id && tag.coords) {
                if (x >= tag.coords.left && x <= tag.coords.left + tag.coords.width
                    && y >= tag.coords.top && y <= tag.coords.top + tag.coords.height) {

                    if (!this.overTag || tag.coords.width < this.overTag.coords.width) {
                        this.overTag = tag;

                        this.overPosition = "inside";

                        if (x < tag.coords.left + tag.coords.width * 0.33
                            && y < tag.coords.top + tag.coords.height * 0.33) {
                            this.overPosition = "before";
                        }

                        if (x > tag.coords.left + tag.coords.width * (1 - 0.33)
                            && y > tag.coords.top + tag.coords.height * (1 - 0.33)) {
                            this.overPosition = "after";
                        }


                    }
                }
            }
            if (tag.body) {
                for (var i = 0; i < tag.body.length; i++) {
                    this.findNextTag(x, y, tag.body[i]);
                }
            }
        }
    }
});

var Tag = JS_CLASS({
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
    }
});

/*простые теги, например a, p, img*/
var SimpleTag = JS_CLASS(Tag, {
    params: {},
    name: "",
    lineOffset: 0,
    constructor: function (param, lineOffset) {
        CP(this, param);
        if(lineOffset)
            this.lineOffset = lineOffset;
        console.log("construct " + this.name, this.lineOffset);
    },
    renderCode: function () {

        var code = this.renderOffset() + this.name + "{";
        var paramArr = [];
        for(var paramName in this.params) {
            if (this.params.hasOwnProperty(paramName)) {
                paramArr.push(paramName + ' = "' + this.params[paramName] + '"');
            }
        }
        code += paramArr.join(", ");
        code += "}\n";
        return code;
    }
});

/* теги со вложенностью, например LiBegin, UList*/
var StructureTag = JS_CLASS(Tag, {
    params: {},
    nameBegin: "",
    nameEnd: "",
    lineOffset: 0,
    constructor: function (param, lineOffset) {
        CP(this, param);
        if(lineOffset)
            this.lineOffset = lineOffset;
    },
    renderCode: function () {
        var code = this.renderOffset() + this.nameBegin + "{";
        var paramArr = [];
        for(var paramName in this.params) {
            if (this.params.hasOwnProperty(paramName)) {
                paramArr.push(paramName + ' = "' + this.params[paramName] + '"');
            }
        }
        code += paramArr.join(", ");
        code += "}\n";

        code += this.renderSubItems();

        code += this.renderOffset() + this.nameEnd + ":\n";
        return code;
    },

    renderSubItems: function (returnType) {
        var code = "";
        if(this.body) {
            for(var i = 0; i < this.body.length; i++) {
                var item = this.body[i];

                var tagObj = null;

                switch(item.name) {
                    case "A":
                        tagObj = new TagA(item, this.lineOffset + 1);
                        break;
                    case "P":
                        tagObj = new TagP(item, this.lineOffset + 1);
                        break;
                    case "Div":
                        tagObj = new TagDiv(item, this.lineOffset + 1);
                        break;
                    case "Image":
                        tagObj = new TagImage(item, this.lineOffset + 1);
                        break;
                    case "Li":
                        tagObj = new TagLi(item, this.lineOffset);
                        break;
                    case "Divs":
                        tagObj = new TagDivs(item, this.lineOffset + 1);
                        break;
                    case "UList":
                        tagObj = new TagUList(item, this.lineOffset + 1);
                        break;
                    case "LiBegin":
                        tagObj = new TagLiBegin(item, this.lineOffset + 1);
                        break;
                }

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

var TagDivs = JS_CLASS(StructureTag, {
    nameBegin: "Divs",
    nameEnd: "DivsEnd",
    renderCode: function () {
        var code = this.renderOffset() + this.nameBegin + "(";

        code += this.nestedClassList.join(", ");
        code += ")\n";

        code += this.renderSubItems();

        code += this.renderOffset() + this.nameEnd + ":\n";
        return code;
    },

    renderHTML: function () {
        var html = '';

        for(var i = 0; i < this.nestedClassList.length; i++) {
            html += '<div tag-id="' + this.id + '" class="'+ this.nestedClassList[i] +'">';
        }

        html += this.renderSubItems('html');

        for(i = 0; i < this.nestedClassList.length; i++) {
            html += '</div>';
        }

        return html;
    }
});

var TagUList = JS_CLASS(StructureTag, {
    nameBegin: "UList",
    nameEnd: "UListEnd",

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

var TagLiBegin = JS_CLASS(StructureTag, {
    nameBegin: "LiBegin",
    nameEnd: "LiEnd",

    renderHTML: function () {
        var html = '<li tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">';

        html += this.renderSubItems('html');

        html += '</li>';
        return html;
    }
});

var TagA = JS_CLASS(SimpleTag, {
    name: "A",

    renderHTML: function () {
        var html = '<a tag-id="' + this.id + '" href="' + this.params.href + '" class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</a>';
        return html;
    }
});

var TagP = JS_CLASS(SimpleTag, {
    name: "P",

    renderHTML: function () {
        var html = '<p tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</p>';
        return html;
    }
});

var TagDiv = JS_CLASS(SimpleTag, {
    name: "Div",
    renderHTML: function () {
        var html = '<div tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</div>';
        return html;
    }
});

var TagImage = JS_CLASS(SimpleTag, {
    name: "Image",
    renderHTML: function () {
        var html = '<img tag-id="' + this.id + '" src="' + (this.params.src ? this.params.src : "") + '" class="' + (this.params.class ? this.params.class : "") + '" alt="' + (this.params.alt ? this.params.alt : "") + '">';
        return html;
    }
});

var TagLi = JS_CLASS(SimpleTag, {
    name: "Li",
    renderHTML: function () {
        var html = '<li tag-id="' + this.id + '" class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</li>';
        return html;
    }
});
