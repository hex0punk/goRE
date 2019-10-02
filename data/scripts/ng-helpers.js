var injector,
    getSelectedComponent,
    applyChanges,
    getService,
    getAllServices,
    getAllServices2;

var inj;

function gorp() {
    if (typeof ng !== "undefined"){
        console.log('App is running Angular 2+')
        setupAngular2Scripts();
    }
    if (typeof angular !== 'undefined'){
        console.log("App is running Angular " + angular.version.full);
        setupAngularScripts();
    }
}


function setupAngular2Scripts(){
    getSelectedComponent = function (){
        if ($0 != null){
            var state = ng.probe($0);
            debugger;
            return state.componentInstance;
        }
        console.log("Select TEST the root element from the inspector");
    };
    applyChanges = function(){
        ng.probe($0).injector.get(ng.coreTokens.ApplicationRef).tick();
    };
    console.log("Available functions:\n- getSelectedComponent()\n- applyChanges()\n");
}

function setupAngularScripts(){
    injector = angular.element(document.body).injector();
    getService = function (serviceName) {
        return injector.get(serviceName)
    };
    applyChanges = function () {
        angular.element(document.body).scope().$apply();
    };
    getAllServices2 = function (mod, r) {
        var inj = angular.element(document).injector().get;
        if (!r) r = {};
        angular.forEach(angular.module(mod).requires, function(m) {getAllServices2(m,r)});
        angular.forEach(angular.module(mod)._invokeQueue, function(a) {
            try { r[a[2][0]] = inj(a[2][0]); } catch (e) {}
        });
        return r;
    };
    getAllServices = function(mod, r) {
        if (!r) {
            r = {};
            if (document.querySelector('[ng-app]'))
                inj = angular.element(document.querySelector('[ng-app]')).injector().get;
            if (document.querySelector('[data-ng-app]'))
                inj = angular.element(document.querySelector('[data-ng-app]')).injector().get;
        }
        angular.forEach(angular.module(mod).requires, function(m) {getAllServices(m,r)});
        angular.forEach(angular.module(mod)._invokeQueue, function(a) {
            try { r[a[2][0]] = inj(a[2][0]); } catch (e) {}
        });
        return r;
    };
    console.log("Available functions:\n- getAllServices()\n- applyChanges()\n- getService()");
    console.log("Available variables:\n- injector");
}
