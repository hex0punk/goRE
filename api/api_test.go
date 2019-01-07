package api

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

const js = `Yb.leaveScope(a,b),b}function ja(a,b){return Yb.beginTimeRange(a,b)}function ka(a){Yb.endTimeRange(a)}function za(a,b){return null}function ya(a){if(0==a._nesting&&!a.hasPendingMicrotasks&&!a.isStable)try{a._nesting++,a.onMicrotaskEmpty.emit(null)}finally{if(a._nesting--,!a.hasPendingMicrotasks)try{a.runOutsideAngular(function(){return a.onStable.emit(null)})}finally{a.isStable=
!0}}}function ua(a){a._inner=a._inner.fork({name:"angular",properties:{isAngularZone:!0},onInvokeTask:function(b,c,e,d,g,k){try{return ra(a),b.invokeTask(e,d,g,k)}finally{a._nesting--,ya(a)}},onInvoke:function(b,c,e,d,g,k,f){try{return ra(a),b.invoke(e,d,g,k,f)}finally{a._nesting--,ya(a)}},onHasTask:function(b,c,e,d){b.hasTask(e,d);c===e&&("microTask"==d.change?(a.hasPendingMicrotasks=d.microTask,ya(a)):"macroTask"==d.change&&(a.hasPendingMacrotasks=d.macroTask))},onHandleError:function(b,c,e,d){return b.handleError(e,
d),a.runOutsideAngular(function(){return a.onError.emit(d)}),!1}})}function ra(a){a._nesting++;a.isStable&&(a.isStable=!1,a.onUnstable.emit(null))}function Da(a){zd=a}function Qa(){if(ib)throw Error("Cannot enable prod mode after platform setup.");ge=!1}function wa(){return ib=!0,ge}function pa(a){if(fa&&!fa.destroyed&&!fa.injector.get(xb,!1))throw Error("There can be only one platform. Destroy the previous one to create a new one.");fa=a.get(sc);a=a.get(Se,null);return a&&a.forEach(function(a){return a()}),
fa}function Ra(a,b,c){void 0===c&&(c=[]);var e=new Ya("Platform: "+b);return function(b){void 0===b&&(b=[]);var d=va();return d&&!d.injector.get(xb,!1)||(a?a(c.concat(b).concat({provide:e,useValue:!0})):pa(bd.resolveAndCreate(c.concat(b).concat({provide:e,useValue:!0})))),Oa(e)}}function Oa(a){var b=va();if(!b)throw Error("No platform exists!");if(!b.injector.get(a,null))throw Error("A platform with a different configuration has been created. Please destroy it first.");return b}function Ta(){fa&&
!fa.destroyed&&fa.destroy()}function va(){return fa&&!fa.destroyed?fa:null}function gb(a,b,c){try{var e=c();return Q(e)?e.catch(function(c){throw b.runOutsideAngular(function(){return a.handleError(c)}),c;}):e}catch(Gc){throw b.runOutsideAngular(function(){return a.handleError(Gc)}),Gc;}}function vb(a,b){b=a.indexOf(b);-1<b&&a.splice(b,1)}function Fb(a,b){var c=cd.get(a);if(c)throw Error("Duplicate module registered for "+a+" - "+c.moduleType.name+" vs "+b.moduleType.name);cd.set(a,b)}function Va(a){var b=`

func TestGetJsFunctionWithHint(t *testing.T){
	enableProdModeFunc := GetJsFunctionWithHint(js, "\"Cannot enable prod mode")
	assert.Equal(t, enableProdModeFunc.Name, "Qa")
	assert.Equal(t, enableProdModeFunc.Raw, `function Qa(){if(ib)throw Error("Cannot enable prod mode after platform setup.");ge=!1}`)
	assert.Equal(t, enableProdModeFunc.Body, `{if(ib)throw Error("Cannot enable prod mode after platform setup.");ge=!1}`)
}