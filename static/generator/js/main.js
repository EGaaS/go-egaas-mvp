var DragAndDrop = JS_CLASS({

    dragging: false,

    constructor: function (param) {
        CP(this, param);

        this.$over = $('<div class="b-over"></div>');
        $(document.body).append(this.$over);

        this.events();
        this.eventsDynamic();



    },

    createDraggableBag: function (html) {
        this.$bag = $('<div class="b-draggable-bag"><div class="b-draggable-bag__inner">'+html+'</div></div>');
        this.$bag.hide();
        $(document.body).append(this.$bag);

    },

    removeDraggableBag: function () {
        if(this.$bag)
            this.$bag.remove();
    },

    moveDraggableBag: function (x, y) {
        this.$bag
            .show()
            .css("top", y + 3)
            .css("left", x + 3);
    },

    events: function () {
        var self = this;

        $(document.body)
            .on("mouseup mouseleave", function (e) {
                if(!self.dragging)
                    return;
                self.dragging = false;
                $(document.body)
                    .removeClass("g-no-select")
                    .removeClass("g-no-overflow-x");
                self.removeDraggableBag();
                self.$container.find(".js-droppable")
                    .removeClass("b-droppable-inside")
                    .removeClass("b-droppable-before")
                    .removeClass("b-droppable-after");

                $(".js-draggable").removeClass("b-draggable_dragging");
            })
            .on("mousemove", function (e) {
                if(!self.dragging)
                    return;
                self.moveDraggableBag(e.pageX, e.pageY);
            });

        // $(".b-over").on("mouseout", function () {
        //     $(this).hide();
        // });

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

                $(".b-over")
                    .show()
                    .css("left", $(this).offset().left)
                    .css("top", $(this).offset().top)
                    .css("width", $(this).outerWidth())
                    .css("height", $(this).outerHeight());
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
                e.stopPropagation();
                if(!self.dragging)
                    return;
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

        $(".js-draggable")
            .off("dragstart")
            .on("dragstart", function (e) {
                e.preventDefault();
            })
            .off("mousedown")
            .on("mousedown", function (e) {
                e.stopPropagation();
                self.dragging = true;
                $(document.body)
                    .addClass("g-no-select")
                    .addClass("g-no-overflow-x");
                self.createDraggableBag(this.outerHTML);
                $(this).addClass("b-draggable_dragging");
            });
    },

    setHTML: function (html) {
        this.html = html;
    },

    render: function () {
        this.$container.html(this.html);
    }
});

$(function() {

    var codeGenerator = new CodeGenerator({
        $container: $(".js-container"),
        $containerWrapper: $(".js-container-wrapper")
    });

    var $input = $(".js-input");
    var $output = $(".js-output");
    var $error = $(".js-error");
    var $codeGenerated = $(".js-code-generated");

    $input.on("change keyup", function() {
        $error.html("");
        try {
            var result = parser.parse($input.val());
            //console.log(result);

            codeGenerator.setJsonData(result);
            $codeGenerated.html(codeGenerator.generateCode());
            codeGenerator.render();
            codeGenerator.eventsDynamic();

            $output.html(JSON.stringify(codeGenerator.json));

        }
        catch (e) {
            $output.html("");
            $error.html(e.message);
        }
    });

    setTimeout(function () {
        $input.trigger("change");
    }, 300);



});



