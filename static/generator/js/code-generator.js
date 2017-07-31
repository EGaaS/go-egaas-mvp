var CodeGenerator = JS_CLASS({

    constructor: function (param) {
        CP(this, param);

    },

    generateCode: function (json) {
        return (new MainTemplate(json)).renderCode();
    },

    generateHTML: function (json) {
        return (new MainTemplate(json)).renderHTML();
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
            html += '<div class="'+ this.nestedClassList[i] +'">';
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
        var html = '<' + tag + ' class="' + (this.params.class ? this.params.class : "") + '">';

        html += this.renderSubItems('html');

        html += '</' + tag + '>';

        return html;
    }
});

var TagLiBegin = JS_CLASS(StructureTag, {
    nameBegin: "LiBegin",
    nameEnd: "LiEnd",

    renderHTML: function () {
        var html = '<li class="' + (this.params.class ? this.params.class : "") + '">';

        html += this.renderSubItems('html');

        html += '</li>';
        return html;
    }
});

var TagA = JS_CLASS(SimpleTag, {
    name: "A",

    renderHTML: function () {
        var html = '<a href="' + this.params.href + '" class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</a>';
        return html;
    }
});

var TagP = JS_CLASS(SimpleTag, {
    name: "P",

    renderHTML: function () {
        var html = '<p class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</p>';
        return html;
    }
});

var TagDiv = JS_CLASS(SimpleTag, {
    name: "Div",
    renderHTML: function () {
        var html = '<div class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</div>';
        return html;
    }
});

var TagImage = JS_CLASS(SimpleTag, {
    name: "Image",
    renderHTML: function () {
        var html = '<img src="' + (this.params.src ? this.params.src : "") + '" class="' + (this.params.class ? this.params.class : "") + '" alt="' + (this.params.alt ? this.params.alt : "") + '">';
        return html;
    }
});

var TagLi = JS_CLASS(SimpleTag, {
    name: "Li",
    renderHTML: function () {
        var html = '<li class="' + (this.params.class ? this.params.class : "") + '">' + this.params.text + '</li>';
        return html;
    }
});

var DragAndDrop = JS_CLASS({
    constructor: function (param) {
        CP(this, param);
        this.events();
    },

    events: function () {
        var self = this;
        this.$container.find(".js-droppable").each(function() {
            $(this).on("mouseover", function (e) {
                e.stopPropagation();
                // self.$container.find(".js-droppable").each(function() {
                //     $(this).removeClass("b-draggable-inside");
                // });
                console.log("over droppable", e);
                //$(e.target).addClass("b-draggable-inside");
                $(this).addClass("b-draggable-inside");
				//$(this).css('backgroundColor', '#a0adff');
            });

            $(this).on("mouseout", function (e) {
                e.stopPropagation();
                // self.$container.find(".js-droppable").each(function() {
                //     $(this).removeClass("b-draggable-inside");
                // });
                console.log("out droppable", e);
                //$(e.target).addClass("b-draggable-inside");
                $(this).removeClass("b-draggable-inside");				
            });
        });


        // this.$container.find(".js-droppable").on("mouseout", function (e) {
        //     e.stopPropagation();
        //     // self.$container.find(".js-droppable").each(function() {
        //     //     $(this).removeClass("b-draggable-inside");
        //     // });
        //     console.log("out to droppable", e);
        //     //$(e.target).removeClass("b-draggable-inside");
        // });
    }
});