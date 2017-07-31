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
                //$(this).addClass("b-draggable-inside-2");
            });

            $(this).on("mouseout", function (e) {
                e.stopPropagation();
                // self.$container.find(".js-droppable").each(function() {
                //     $(this).removeClass("b-draggable-inside");
                // });
                console.log("out droppable", e);
                //$(e.target).addClass("b-draggable-inside");
                $(this).removeClass("b-draggable-inside");
                $(this).removeClass("b-draggable-before");
                $(this).removeClass("b-draggable-after");
            });

            $(this).on("mousemove", function (e) {
                e.stopPropagation();
                // self.$container.find(".js-droppable").each(function() {
                //     $(this).removeClass("b-draggable-inside");
                // });
                console.log("move droppable", e);

                var width = $(this).outerWidth();
                var height = $(this).outerHeight();

                //var offsetX = e.offsetX - this.offsetLeft;
                //var offsetY = e.offsetY - this.offsetTop;

                var offsetX = e.pageX - $(this).offset().left;
                var offsetY = e.pageY - $(this).offset().top;

                // inside a scrolling pane:
                //
                // var x = (evt.pageX - $('#element').offset().left) + self.frame.scrollLeft();
                // var y = (evt.pageY - $('#element').offset().top) + self.frame.scrollTop();

                console.log("el width", width, height, offsetX, offsetY);

                if(offsetX < width * 0.33 && offsetY < height * 0.33) {
                    $(this).addClass("b-draggable-before");
                    $(this).removeClass("b-draggable-inside");
                    $(this).removeClass("b-draggable-after");
                }
                else {
                    if (offsetX > width * (1 - 0.33) && offsetY > height * (1 - 0.33)) {
                        $(this).addClass("b-draggable-after");
                        $(this).removeClass("b-draggable-inside");
                        $(this).removeClass("b-draggable-before");
                    }
                    else {
                        $(this).addClass("b-draggable-inside");
                        $(this).removeClass("b-draggable-before");
                        $(this).removeClass("b-draggable-after");
                    }
                }

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

$(function() {
    //alert("main");

    $(".js-source-element").draggable({
        //addClasses: false,
        cursor: "move",
        cursorAt: { left: -1, top: -1 },
        //refreshPositions: true,
        insertBeforeAndAfter: true,
        helper: function () {
            return $("<div style='padding: 10px; border: 1px dashed #CCC; white-space: nowrap; background: rgba(255, 255, 255, 0.7); cursor: move;'>Тащим направо в контейнер</div>");
        },
        over: function (event, ui) {
            //console.log("event", event);
            //console.log("ui", ui);
        }

    });

    $(".js-tag").draggable({
        //addClasses: false,
        cursor: "move",
        helper: "clone",
        //refreshPositions: true,
        over: function (event, ui) {
            console.log("event", event);
            console.log("ui", ui);
        }

    });

    $(".js-tag").droppable({
        greedy: true,
        tolerance: "pointer",
        //tolerance: "touch",
        insertBeforeAndAfter: true,
        accept: function() {
            //console.log(this.tagName);
            switch(this.tagName) {
                case "DIV":
                    return true;
            };
            return false;
        },
        acceptBefore: '*',
        acceptAfter: '*',
        classes: {
            "ui-droppable-hover": 'b-draggable-inside',
            "ui-droppable-before": 'b-draggable-before',
            "ui-droppable-after": 'b-draggable-after'
        },
        //hoverClass: 'b-draggable-inside',
        over: function (event, ui) {

            //console.log(this); //the 'this' under over event
            //console.log(event);
        }
    });

    $(".js-container-").droppable({
        greedy: true,
        tolerance: "pointer",
        //tolerance: "touch",
        insertBeforeAndAfter: true,
        accept: function() {
            return true;
        },
        acceptBefore: function() {
            return false;
        },
        acceptAfter: function() {
            return false;
        },
        classes: {
            "ui-droppable-hover": 'b-draggable-inside',
            "ui-droppable-before": 'b-draggable-before',
            "ui-droppable-after": 'b-draggable-after'
        },
        //hoverClass: 'b-draggable-inside',
        over: function (event, ui) {

            //console.log(this); //the 'this' under over event
            //console.log(event);
        }
    });


    $(".js-container").on("mousemove", ".js-tag", function (e) {
        //console.log("move", e);
    });

    var dragAndDrop = new DragAndDrop({
        $container: $(".js-container2")
    });


});



