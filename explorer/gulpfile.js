'use strict';

var gulp = require('gulp'),
    changed = require('gulp-changed'),
    sass = require('gulp-sass'),
    csso = require('gulp-csso'),
    autoprefixer = require('gulp-autoprefixer'),
    browserify = require('browserify'),
    watchify = require('watchify'),
    source = require('vinyl-source-stream'),
    buffer = require('vinyl-buffer'),
    reactify = require('reactify'),
    uglify = require('gulp-uglify'),
    del = require('del'),
    notify = require('gulp-notify'),
    browserSync = require('browser-sync'),
    reload = browserSync.reload,

    WatchNetwork = require("gulp-watch-network"),

    p = {
      jsx: './src/app.jsx',
      scss: 'styles/main.scss',
      bundle: 'app.js',
      distJs: 'dist/js',
      distCss: 'dist/css'
    };

gulp.task('clean', function(cb) {
  del(['dist'], cb);
});

gulp.task('browserSync', function() {
  browserSync({
    server: {
      baseDir: './'
    }
  })
});

gulp.task('watchify', function() {
  var args = watchify.args
  args.paths = ['./node_modules','./src']

  var bundler = watchify(browserify(p.jsx, args));

  function rebundle() {
    return bundler
      .bundle()
      .on('error', notify.onError())
      .pipe(source(p.bundle))
      .pipe(gulp.dest(p.distJs))
      .pipe(reload({stream: true}));
  }

  bundler.transform(reactify).on('update', rebundle);

  return rebundle();
});

// SALUT
gulp.task('watch-network', function() {
  var watch = WatchNetwork({
    gulp: gulp,
    host: '192.168.30.1',
    port: '4000',
    configs: [{
      tasks: 'build',
      onLoad: true
    }, {
      patterns: ['src/**.js*'],
      tasks: 'browserify'
    }, {
      patterns: 'styles/*.scss',
      tasks: ['styles', 'fonts']
    }]
  })

  watch.initialize();
});

gulp.task('browserify', function() {
  browserify(p.jsx, {
      basedir: __dirname,
      debug: true,
      paths: ['./node_modules','./src']
    })
    .transform(reactify)
    .bundle()
    .pipe(source(p.bundle))
    .pipe(buffer())
    //.pipe(uglify())
    .pipe(gulp.dest(p.distJs));
});

gulp.task('styles', function() {
  return gulp.src(p.scss)
    .pipe(changed(p.distCss))
    .pipe(sass({errLogToConsole: true}))
    .on('error', notify.onError())
    .pipe(autoprefixer('last 1 version'))
    //.pipe(csso())
    .pipe(gulp.dest(p.distCss))
    .pipe(reload({stream: true}));
});

gulp.task('watchTask', function() {
  gulp.watch(p.scss, ['styles']);
});

gulp.task('fonts', function() { 
    return gulp.src(['./bower_components/bootstrap-sass/fonts/**.*', './bower_components/fontawesome/fonts/**.*']) 
        .pipe(gulp.dest('./dist/fonts')); 
});

gulp.task('watch', ['clean'], function() {
  gulp.start(['browserSync', 'watchTask', 'watchify', 'styles', 'fonts']);
});

gulp.task('build', ['clean'], function() {
  process.env.NODE_ENV = 'development';
  gulp.start(['browserify', 'styles', 'fonts']);

  gulp.src('src/index.html').pipe(gulp.dest('dist'));
});

gulp.task('default', function() {
  console.log('Run "gulp watch or gulp build"');
});
