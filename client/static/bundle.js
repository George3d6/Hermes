/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// identity function for calling harmony imports with the correct context
/******/ 	__webpack_require__.i = function(value) { return value; };
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, {
/******/ 				configurable: false,
/******/ 				enumerable: true,
/******/ 				get: getter
/******/ 			});
/******/ 		}
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = 8);
/******/ })
/************************************************************************/
/******/ ([
/* 0 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
exports.displayError = function (message) {
    document.getElementById('error_holder').innerHTML = '';
    document.getElementById('error_holder').insertAdjacentHTML('afterbegin', "\n    <div class=\"ui error message\">\n    <i class=\"close icon\" id=\"close_message\"></i>\n        <div class=\"header\">\n            There was an error\n        </div>\n        <p>\n            " + message + "\n        </p>\n    </div>\n    ");
    setTimeout(function () {
        document.getElementById('error_holder').innerHTML = '';
    }, 2600);
};
exports.displayMessage = function (message) {
    document.getElementById('error_holder').innerHTML = '';
    document.getElementById('error_holder').insertAdjacentHTML('afterbegin', "\n    <div class=\"ui info message\">\n    <i class=\"close icon\" id=\"close_message\"></i>\n        <div class=\"header\">\n            Success\n        </div>\n        <p>\n            " + message + "\n        </p>\n    </div>\n    ");
    setTimeout(function () {
        document.getElementById('error_holder').innerHTML = '';
    }, 2600);
};

/***/ }),
/* 1 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
var uploadFile_1 = __webpack_require__(5);
var listFiles_1 = __webpack_require__(3);
var router_1 = __webpack_require__(4);
var file_1 = __webpack_require__(2);
var authenticate_1 = __webpack_require__(6);
var views_1 = __webpack_require__(0);
uploadFile_1["default"]('upload_form');
authenticate_1["default"]();
var initFileListView = function () {
    listFiles_1["default"]().then(function (answer) {
        document.getElementById('file_list').innerHTML = '';
        var files = [];
        answer.split('#|#').forEach(function (stringifiedFile) {
            if (stringifiedFile === '') {
                return;
            }
            var fileElements = stringifiedFile.split('|#|');
            files.push(new file_1["default"](fileElements[0], fileElements[1], parseInt(fileElements[2])));
            return files;
        });
        document.getElementById('file_list').insertAdjacentHTML('beforeend', "<p id=\"close_file_view_holder\"><a href=\"/#!/empty/\" id=\"close_file_view\"><i class=\"close icon big\"></i></a></p>");
        files.forEach(function (f) {
            var fileDies = new Date(f.death * 1000);
            var fileDiesString = fileDies.getUTCFullYear() + "-" + fileDies.getUTCMonth() + "-" + fileDies.getUTCDay() + " " + fileDies.getUTCHours() + ":" + fileDies.getUTCMinutes() + ":" + fileDies.getUTCSeconds();
            document.getElementById('file_list').insertAdjacentHTML('beforeend', "\n            <div class=\"item file_in_list\">\n             <i class=\"huge cloud download middle aligned icon cursor_hover files_icon\" id=\"download_" + f.name + "\" onclick=\"window.location='/get/file/?file=" + f.name + "';\"></i>\n             <i class=\"huge trash middle aligned icon cursor_hover trashed_icon\" id=\"trash_" + f.name + "\"></i>\n              <div class=\"content\" id=\"" + f.name + "\" class=\"inline_content\">\n                <p>\n                    " + f.name + "\n                </p>\n                <p>\n                    Valid until: " + fileDiesString + "\n                </p>\n                <div>\n                    Compression: " + f.compression + "\n                </div>\n                </div>\n            </div>\n            ");
            document.getElementById("trash_" + f.name).addEventListener('click', function (e) {
                var xhr = new XMLHttpRequest();
                xhr.open('GET', '/delete/file/?file=' + f.name);
                xhr.send(null);
                xhr.onload = function () {
                    console.log(xhr.responseText);
                    var res = JSON.parse(xhr.responseText);
                    if (res.status === "error") {
                        views_1.displayError(res.message);
                    } else {
                        views_1.displayMessage(res.message);
                    }
                };
                console.log("Working");
                initFileListView();
            });
        });
    })["catch"](function (err) {
        console.log("Something went horribly wrong: " + err);
    });
};
window['initFileListView'] = initFileListView;
router_1["default"].on('/files/', function () {
    initFileListView();
}).resolve();
router_1["default"].on('/', function () {
    initFileListView();
}).resolve();
router_1["default"].on('/empty/', function () {
    document.getElementById('file_list').innerHTML = '';
    console.log("HERE:", document.getElementById('file_list').innerHTML);
}).resolve();

/***/ }),
/* 2 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
var FileModel = function () {
    function FileModel(name, compression, death, size) {
        if (size === void 0) {
            size = 0;
        }
        this.name = name;
        this.death = death;
        this.compression = compression;
        this.size = size;
    }
    return FileModel;
}();
exports["default"] = FileModel;

/***/ }),
/* 3 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
var views_1 = __webpack_require__(0);
var listFiles = function () {
    return new Promise(function (resolve, reject) {
        var xhr = new XMLHttpRequest();
        xhr.open("GET", "/get/list/");
        xhr.onload = function () {
            try {
                var res = JSON.parse(xhr.responseText);
                if (res.status === "error") {
                    views_1.displayError(res.message);
                } else {
                    views_1.displayMessage(res.message);
                }
                reject(Error('Could not load files:' + xhr.statusText));
            } catch (e) {
                resolve(xhr.responseText);
            }
        };
        xhr.send(null);
    });
};
exports["default"] = listFiles;

/***/ }),
/* 4 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
var navigo_1 = __webpack_require__(7);
var root = null;
var useHash = true;
var hash = '#!';
var router = new navigo_1["default"](root, useHash, hash);
exports["default"] = router;

/***/ }),
/* 5 */
/***/ (function(module, exports, __webpack_require__) {

"use strict";


exports.__esModule = true;
//Uploads a file and some metadata based on a pre-defined form
var uploadFile = function (form_id) {
    if (window.navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
        var uploadForm_1 = document.getElementById(form_id);
        document.getElementById("submit_form").addEventListener("click", function (e) {
            e.preventDefault();
            var reader = new FileReader();
            reader.readAsArrayBuffer(document.getElementById('file').files[0]);
            reader.onload = function (evt) {
                var formData = new FormData(uploadForm_1);
                var isPublic = String(document.getElementById('public_switch').checked);
                formData.append('file', evt.target.result);
                formData.append('compression', document.getElementById('compression').value);
                formData.append('ispublis', isPublic);
                formData.append('isAsync', "true");
                alert(' the form value is:  ' + formData.get('ispublis'));
                var xhr = new XMLHttpRequest();
                xhr.open("POST", "/post/file/");
                xhr.send(formData);
                xhr.onreadystatechange = function () {
                    console.log(xhr.responseText + '  \n status is: ' + xhr.statusText);
                };
            };
        });
    } else {
        console.log("Warnning, your browser does not support asynchronous upload of large files");
    }
};
exports["default"] = uploadFile;

/***/ }),
/* 6 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__views__ = __webpack_require__(0);
/* harmony import */ var __WEBPACK_IMPORTED_MODULE_0__views___default = __webpack_require__.n(__WEBPACK_IMPORTED_MODULE_0__views__);


const enableAuthenticationForm = () => {

    document.getElementById("submit_auth").addEventListener("click", e => {
        e.preventDefault();
        const identifier = document.getElementById("identifier_field").value;
        const credentials = document.getElementById("credentials_field").value;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/get/authentication/?identifier=${identifier}&credentials=${credentials}`);
        xhr.send();
        xhr.onreadystatechange = () => {
            $('#sign_in_form_modal').modal('hide');
            const res = JSON.parse(xhr.responseText);
            if (res.status === "error") {
                __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__views__["displayError"])(res.message);
            } else {
                __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__views__["displayMessage"])(res.message);
            }
        };
    });

    document.getElementById("submit_auth_make").addEventListener("click", e => {
        e.preventDefault();
        const identifier = document.getElementById("identifier_field_make").value;
        const credentials = document.getElementById("credentials_field_make").value;
        const uploadNumber = document.getElementById("uploadNumber_field_make").value;
        const uploadSize = document.getElementById("uploadSize_field_make").value;
        const reader = document.getElementById("reader_field_make").checked;
        const writer = document.getElementById("writer_field_make").checked;
        const admin = document.getElementById("admin_field_make").checked;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/post/token/?identifier=${identifier}&credentials=${credentials}&uploadNumber=${uploadNumber}` + `&uploadSize=${uploadSize}&reader=${reader}&writer=${writer}&admin=${admin}`);
        xhr.send();
        xhr.onreadystatechange = () => {
            $('#permission_view').modal('hide');
            if (res.status === "error") {
                __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__views__["displayError"])(res.message);
            } else {
                __webpack_require__.i(__WEBPACK_IMPORTED_MODULE_0__views__["displayMessage"])(res.message);
            }
        };
    });

    document.getElementById("close_permission_view_holder").addEventListener("click", e => {
        $('#permission_view').modal('hide');
    });
};

/* harmony default export */ __webpack_exports__["default"] = (enableAuthenticationForm);

/***/ }),
/* 7 */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
Object.defineProperty(__webpack_exports__, "__esModule", { value: true });
//Here because there is a bug with webpack and/or typescript and importing from node_modules
function isPushStateAvailable() {
  return !!(typeof window !== 'undefined' && window.history && window.history.pushState);
}

function Navigo(r, useHash, hash) {
  this.root = null;
  this._routes = [];
  this._useHash = useHash;
  this._hash = typeof hash === 'undefined' ? '#' : hash;
  this._paused = false;
  this._destroyed = false;
  this._lastRouteResolved = null;
  this._notFoundHandler = null;
  this._defaultHandler = null;
  this._usePushState = !useHash && isPushStateAvailable();
  this._onLocationChange = this._onLocationChange.bind(this);

  if (r) {
    this.root = useHash ? r.replace(/\/$/, '/' + this._hash) : r.replace(/\/$/, '');
  } else if (useHash) {
    this.root = this._cLoc().split(this._hash)[0].replace(/\/$/, '/' + this._hash);
  }

  this._listen();
  this.updatePageLinks();
}

function clean(s) {
  if (s instanceof RegExp) return s;
  return s.replace(/\/+$/, '').replace(/^\/+/, '/');
}

function regExpResultToParams(match, names) {
  if (names.length === 0) return null;
  if (!match) return null;
  return match.slice(1, match.length).reduce((params, value, index) => {
    if (params === null) params = {};
    params[names[index]] = value;
    return params;
  }, null);
}

function replaceDynamicURLParts(route) {
  var paramNames = [],
      regexp;

  if (route instanceof RegExp) {
    regexp = route;
  } else {
    regexp = new RegExp(clean(route).replace(Navigo.PARAMETER_REGEXP, function (full, dots, name) {
      paramNames.push(name);
      return Navigo.REPLACE_VARIABLE_REGEXP;
    }).replace(Navigo.WILDCARD_REGEXP, Navigo.REPLACE_WILDCARD) + Navigo.FOLLOWED_BY_SLASH_REGEXP, Navigo.MATCH_REGEXP_FLAGS);
  }
  return { regexp, paramNames };
}

function getUrlDepth(url) {
  return url.replace(/\/$/, '').split('/').length;
}

function compareUrlDepth(urlA, urlB) {
  return getUrlDepth(urlB) - getUrlDepth(urlA);
}

function findMatchedRoutes(url, routes = []) {
  return routes.map(route => {
    var _replaceDynamicURLPar = replaceDynamicURLParts(route.route),
        regexp = _replaceDynamicURLPar.regexp,
        paramNames = _replaceDynamicURLPar.paramNames;

    var match = url.match(regexp);
    var params = regExpResultToParams(match, paramNames);

    return match ? { match, route, params } : false;
  }).filter(m => m);
}

function match(url, routes) {
  return findMatchedRoutes(url, routes)[0] || false;
}

function root(url, routes) {
  var matched = findMatchedRoutes(url, routes.filter(route => {
    let u = clean(route.route);

    return u !== '' && u !== '*';
  }));
  var fallbackURL = clean(url);

  if (matched.length > 0) {
    return matched.map(m => clean(url.substr(0, m.match.index))).reduce((root, current) => {
      return current.length < root.length ? current : root;
    }, fallbackURL);
  }
  return fallbackURL;
}

function isHashChangeAPIAvailable() {
  return !!(typeof window !== 'undefined' && 'onhashchange' in window);
}

function extractGETParameters(url) {
  return url.split(/\?(.*)?$/).slice(1).join('');
}

function getOnlyURL(url, useHash, hash) {
  var onlyURL = url.split(/\?(.*)?$/)[0];

  if (typeof hash === 'undefined') {
    // To preserve BC
    hash = '#';
  }

  if (isPushStateAvailable() && !useHash) {
    onlyURL = onlyURL.split(hash)[0];
  }

  return onlyURL;
}

function manageHooks(handler, route, params) {
  if (route && route.hooks && typeof route.hooks === 'object') {
    if (route.hooks.before) {
      route.hooks.before((shouldRoute = true) => {
        if (!shouldRoute) return;
        handler();
        route.hooks.after && route.hooks.after(params);
      }, params);
    } else if (route.hooks.after) {
      handler();
      route.hooks.after && route.hooks.after(params);
    }
    return;
  }
  handler();
};

function isHashedRoot(url, useHash, hash) {
  if (isPushStateAvailable() && !useHash) {
    return false;
  }

  if (!url.match(hash)) {
    return false;
  }

  let split = url.split(hash);

  if (split.length < 2 || split[1] === '') {
    return true;
  }

  return false;
};

Navigo.prototype = {
  helpers: {
    match,
    root,
    clean
  },
  navigate: function (path, absolute) {
    var to;

    path = path || '';
    if (this._usePushState) {
      to = (!absolute ? this._getRoot() + '/' : '') + path.replace(/^\/+/, '/');
      to = to.replace(/([^:])(\/{2,})/g, '$1/');
      history[this._paused ? 'replaceState' : 'pushState']({}, '', to);
      this.resolve();
    } else if (typeof window !== 'undefined') {
      path = path.replace(new RegExp('^' + this._hash), '');
      window.location.href = window.location.href.replace(/#$/, '').replace(new RegExp(this._hash + '.*$'), '') + this._hash + path;
    }
    return this;
  },
  on: function (...args) {
    if (typeof args[0] === 'function') {
      this._defaultHandler = { handler: args[0], hooks: args[1] };
    } else if (args.length >= 2) {
      if (args[0] === '/') {
        let func = args[1];

        if (typeof args[1] === 'object') {
          func = args[1].uses;
        }

        this._defaultHandler = { handler: func, hooks: args[2] };
      } else {
        this._add(args[0], args[1], args[2]);
      }
    } else if (typeof args[0] === 'object') {
      let orderedRoutes = Object.keys(args[0]).sort(compareUrlDepth);

      orderedRoutes.forEach(route => {
        this.on(route, args[0][route]);
      });
    }
    return this;
  },
  off: function (handler) {
    if (this._defaultHandler !== null && handler === this._defaultHandler.handler) {
      this._defaultHandler = null;
    } else if (this._notFoundHandler !== null && handler === this._notFoundHandler.handler) {
      this._notFoundHandler = null;
    }
    this._routes = this._routes.reduce((result, r) => {
      if (r.handler !== handler) result.push(r);
      return result;
    }, []);
    return this;
  },
  notFound: function (handler, hooks) {
    this._notFoundHandler = { handler, hooks: hooks };
    return this;
  },
  resolve: function (current) {
    var handler, m;
    var url = (current || this._cLoc()).replace(this._getRoot(), '');

    if (this._useHash) {
      url = url.replace(new RegExp('^\/' + this._hash), '/');
    }

    let GETParameters = extractGETParameters(current || this._cLoc());
    let onlyURL = getOnlyURL(url, this._useHash, this._hash);

    if (this._paused || this._lastRouteResolved && onlyURL === this._lastRouteResolved.url && GETParameters === this._lastRouteResolved.query) {
      return false;
    }

    m = match(onlyURL, this._routes);

    if (m) {
      this._lastRouteResolved = { url: onlyURL, query: GETParameters };
      handler = m.route.handler;
      manageHooks(() => {
        m.route.route instanceof RegExp ? handler(...m.match.slice(1, m.match.length)) : handler(m.params, GETParameters);
      }, m.route, m.params);
      return m;
    } else if (this._defaultHandler && (onlyURL === '' || onlyURL === '/' || onlyURL === this._hash || isHashedRoot(onlyURL, this._useHash, this._hash))) {
      manageHooks(() => {
        this._lastRouteResolved = { url: onlyURL, query: GETParameters };
        this._defaultHandler.handler(GETParameters);
      }, this._defaultHandler);
      return true;
    } else if (this._notFoundHandler) {
      manageHooks(() => {
        this._lastRouteResolved = { url: onlyURL, query: GETParameters };
        this._notFoundHandler.handler(GETParameters);
      }, this._notFoundHandler);
    }
    return false;
  },
  destroy: function () {
    this._routes = [];
    this._destroyed = true;
    clearTimeout(this._listenningInterval);
    if (typeof window !== 'undefined') {
      window.removeEventListener('popstate', this._onLocationChange);
      window.removeEventListener('hashchange', this._onLocationChange);
    }
  },
  updatePageLinks: function () {
    var self = this;

    if (typeof document === 'undefined') return;

    this._findLinks().forEach(link => {
      if (!link.hasListenerAttached) {
        link.addEventListener('click', function (e) {
          var location = self.getLinkPath(link);

          if (!self._destroyed) {
            e.preventDefault();
            self.navigate(clean(location));
          }
        });
        link.hasListenerAttached = true;
      }
    });
  },
  generate: function (name, data = {}) {
    var result = this._routes.reduce((result, route) => {
      var key;

      if (route.name === name) {
        result = route.route;
        for (key in data) {
          result = result.replace(':' + key, data[key]);
        }
      }
      return result;
    }, '');

    return this._useHash ? this._hash + result : result;
  },
  link: function (path) {
    return this._getRoot() + path;
  },
  pause: function (status = true) {
    this._paused = status;
  },
  resume: function () {
    this.pause(false);
  },
  disableIfAPINotAvailable: function () {
    if (!isPushStateAvailable()) {
      this.destroy();
    }
  },
  lastRouteResolved() {
    return this._lastRouteResolved;
  },
  getLinkPath(link) {
    return link.pathname || link.getAttribute('href');
  },
  _add: function (route, handler = null, hooks = null) {
    if (typeof route === 'string') {
      route = encodeURI(route);
    }
    if (typeof handler === 'object') {
      this._routes.push({
        route,
        handler: handler.uses,
        name: handler.as,
        hooks: hooks || handler.hooks
      });
    } else {
      this._routes.push({ route, handler, hooks: hooks });
    }
    return this._add;
  },
  _getRoot: function () {
    if (this.root !== null) return this.root;
    this.root = root(this._cLoc().split('?')[0], this._routes);
    return this.root;
  },
  _listen: function () {
    if (this._usePushState) {
      window.addEventListener('popstate', this._onLocationChange);
    } else if (isHashChangeAPIAvailable()) {
      window.addEventListener('hashchange', this._onLocationChange);
    } else {
      let cached = this._cLoc(),
          current,
          check;

      check = () => {
        current = this._cLoc();
        if (cached !== current) {
          cached = current;
          this.resolve();
        }
        this._listenningInterval = setTimeout(check, 200);
      };
      check();
    }
  },
  _cLoc: function () {
    if (typeof window !== 'undefined') {
      if (typeof window.__NAVIGO_WINDOW_LOCATION_MOCK__ !== 'undefined') {
        return window.__NAVIGO_WINDOW_LOCATION_MOCK__;
      }
      return clean(window.location.href);
    }
    return '';
  },
  _findLinks: function () {
    return [].slice.call(document.querySelectorAll('[data-navigo]'));
  },
  _onLocationChange: function () {
    this.resolve();
  }
};

Navigo.PARAMETER_REGEXP = /([:*])(\w+)/g;
Navigo.WILDCARD_REGEXP = /\*/g;
Navigo.REPLACE_VARIABLE_REGEXP = '([^\/]+)';
Navigo.REPLACE_WILDCARD = '(?:.*)';
Navigo.FOLLOWED_BY_SLASH_REGEXP = '(?:\/$|$)';
Navigo.MATCH_REGEXP_FLAGS = '';

/* harmony default export */ __webpack_exports__["default"] = (Navigo);

/***/ }),
/* 8 */
/***/ (function(module, exports, __webpack_require__) {

module.exports = __webpack_require__(1);


/***/ })
/******/ ]);