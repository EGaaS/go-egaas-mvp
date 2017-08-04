
$(function() {

    var $input = $(".js-input");
    var $output = $(".js-output");
    var $error = $(".js-error");
    var $codeGenerated = $(".js-code-generated");

    var codeGenerator = new CodeGenerator.Controller({
        $container: $(".js-container"),
        $containerWrapper: $(".js-container-wrapper"),
        $codeGenerated: $codeGenerated,
        $output: $output
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



});



