//if(typeof console == 'undefined') console = { "log": function() { } };

/*if(typeof console == "undefined")
 var console = {
 log: function(p) {
 $(".js-console").prepend("<div>"+p+"</div>");
 //$(".js-console").scrollTop(100000);
 }
 };
 */



// var log = function(p) {
//     if(typeof p == 'object')
//         p = JSON.stringify(p);
//     $(".js-console").prepend("<div>"+p+"</div>");
//     //$(".js-console").scrollTop(100000);
// }

if(!Function.prototype.bind)
    Function.prototype.bind=function(scope)
    {
        var fn=this;
        if(arguments.length<2){
            return function(){
                return fn.apply(scope, arguments);
            };
        }else{
            var args=Array.prototype.slice.call(arguments,1);
            return function(){
                return fn.apply(scope, args.concat(Array.prototype.slice.call(arguments,0)));
            };
        }
    };

function CP(dst,src,def)
{
    var i;
    for(i in def)
        if(def.hasOwnProperty(i))
            dst[i]=def[i];
    for(i in src)
        if(src.hasOwnProperty(i))
            dst[i]=src[i];
    return dst;
}

// копирует только cтроки и числа.
function CP_SIMPLE(dst,src)
{
    var i, o;
    for(i in src) {
        if(src.hasOwnProperty(i)) {
            o = src[i];
            if ( typeof(o) == 'string' || typeof(o) == 'number') dst[i]=o;
        }
    }
    return dst;
}

function TPL(tpl,dict) {
    if(typeof tpl == "undefined")
        return '';
    return tpl.replace(
        /\{([\w-]+)\}/g,
        function(str,key){
            return dict[key]||(typeof(dict[key])=='number'?0:'');
        }
    );
}

function SUPER(scope,args)
{
    var fn=args.callee.SUPER;
    if(arguments.length>2)
        args=Array.prototype.slice.call(arguments,2);
    return fn && fn.apply(scope,args);
}


function JS_CLASS(superclass, overrides)
{
	var subclass, F, supp, subp, i, s, d;

	// normalize params
	i=superclass instanceof Function;
	overrides=overrides || !i && superclass || {};
	superclass=i && superclass || null;

	// find subclass constructor
	if(overrides.constructor!=Object.prototype.constructor)
		subclass=overrides.constructor, delete overrides.constructor;
	else if(superclass)
		subclass=function(){ arguments.callee.SUPER.apply(this,arguments); };
	else
		subclass=function(){};

	if(!superclass){
		subclass.prototype=overrides;
	}else{
		// prototyping
		F=function(){};
		supp=F.prototype=superclass.prototype;
		subp=subclass.prototype=new F();
		subp.constructor=subclass;
		subclass.SUPER=superclass;

		if(supp.constructor==Object.prototype.constructor)
			supp.constructor=superclass;

		for(i in overrides){
			s=overrides[i], d=subp[i];
			if(s instanceof Function && d instanceof Function)
				s.SUPER=d;
			subp[i]=s;
		}
	}

	return this.constructor==arguments.callee ? new subclass() : subclass;
}
