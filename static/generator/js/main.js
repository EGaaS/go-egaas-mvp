var DragAndDrop = JS_CLASS({

    dragging: false,

    constructor: function (param) {
        CP(this, param);
        this.events();
    },

    createDraggableBag: function (html) {
        this.$bag = $('<div class="b-draggable-bag"><div class="b-draggable-bag__inner">'+html+'</div></div>');
        $(document.body).append(this.$bag);
    },

    removeDraggableBag: function () {
        this.$bag.remove();
    },

    moveDraggableBag: function (x, y) {
        this.$bag
            .css("top", y + 3)
            .css("left", x + 3);
    },

    events: function () {
        var self = this;

        this.$container.find(".js-droppable")
            .off("mouseout")
            .on("mouseout", function (e) {
                e.stopPropagation();
                if(!self.dragging)
                    return;
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
                            .addClass("b-draggable-after")
                            .removeClass("b-draggable-inside")
                            .removeClass("b-draggable-before");
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

        $(document.body)
            .on("mouseup mouseleave", function (e) {
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

    }
});

$(function() {
    var dragAndDrop = new DragAndDrop({
        $container: $(".js-container")
    });
});



