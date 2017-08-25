
$(function() {

    var $input = $(".js-input");
    var $output = $(".js-output");
    var $error = $(".js-error");
    var $codeGenerated = $(".js-code-generated");

    var $visualEditorOpen = $(".js-visual-editor-open");
    var $visualEditorCancel = $(".js-visual-editor-cancel");
    var $visualEditorSave = $(".js-visual-editor-save");

    var codeGenerator = new CodeGenerator.Controller({
        $container: $(".js-container"),
        $containerWrapper: $(".js-container-wrapper"),
        $codeGenerated: $codeGenerated,
        $output: $output,
        $instrumentPanel: $(".js-instrument-panel")
    });

    setInterval(function () {
        codeGenerator.printJSON();
    }, 1000);

    $input.on("change keyup", function() {
        $error.html("");
        try {
            var result = parser.parse($input.val());
            //console.log(result);

            codeGenerator.setJsonData(result);
            codeGenerator.generateCode();
            codeGenerator.render();

        }
        catch (e) {
            $output.html("");
            $error.html(e.message);
        }
    });



    setTimeout(function () {
        $input.trigger("change");
    }, 300);

    $visualEditorOpen.on("click", function () {
        $(".js-visual-editor-on").show();
        $(".js-visual-editor-off").hide();
        $error.html("");

        try {
            var result = parser.parse(editor.getValue());
            //console.log(result);
            codeGenerator.setJsonData(result);
            codeGenerator.generateCode();
            codeGenerator.render();
        }
        catch (e) {
            $output.html("");
            $error.html(e.message);
        }
    });

    $visualEditorCancel.on("click", function () {
        $(".js-visual-editor-off").show();
        $(".js-visual-editor-on").hide();
    });

    $visualEditorSave.on("click", function () {
        $(".js-visual-editor-off").show();
        $(".js-visual-editor-on").hide();
        editor.setValue(codeGenerator.getCode());
    });


});



