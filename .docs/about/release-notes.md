# Release Notes

---

## Upgrading

To upgrade Polly to the latest version, use `curl install`:

    curl -sSL https://dl.bintray.com/emccode/rexray/install | sh

Use `polly version` to determine the currently installed version of Polly:

    $ polly version
    Polly
    ----------
    Binary: /usr/bin/polly
    SemVer: 0.1.1-rc1
    OsArch: Linux-x86_64
    Branch: master
    Commit: f34b81a0dee2b22ff6cc4ec10fadd2246597fe22
    Formed: Thu, 30 Jun 2016 18:09:25 UTC

    libStorage
    ----------
    SemVer: ..
    OsArch: Linux-x86_64
    Branch: (detached from ca8ecf1
    Commit: ca8ecf16a73fff0be5918deffa42c8f2983b41c2
    Formed: Thu, 30 Jun 2016 18:08:40 UTC

## Version 0.1.1 (2016/07/12)
Polly 0.1.1 introduces a simple and easy to stand up Vagrant development, test
and demonstration environment for VirtualBox to exercise operations between REX-Ray
acting as a front-end for Polly which is acting as a central management point for
all volume related operations.

To obtain more information and get directions for starting and operating this
Vagrant environment, please checkout the
[ReadTheDocs](http://polly-scheduler.readthedocs.io/en/latest/) for Polly.

### libstorage update
This release also updates the reference to link against libstorage 0.1.5. This fixes a performance issue that was discovered during testing of REX-Ray and Polly. For more information, you can visit the libstorage github page [here](https://github.com/emccode/libstorage).

### New Features
* Vagrant up demo ([#31](https://github.com/emccode/polly/issues/31), [#112](https://github.com/emccode/polly/issues/112)

### Enhancements
* Recommendations for configuration file format ([#106](https://github.com/emccode/polly/issues/106))
* Project dependencies updated ([#102](https://github.com/emccode/polly/issues/102))

### Bug Fixes
* Helpful error messages are lost on certain operations ([#88](https://github.com/emccode/polly/issues/88)),
([#82](https://github.com/emccode/polly/issues/82))
* Fixed dead links in documentation ([#85](https://github.com/emccode/polly/issues/85))
* Fixed Makefile use for uname ([#84](https://github.com/emccode/polly/issues/84))
* Fixed panics to be written to the log ([#68](https://github.com/emccode/polly/issues/68))
* Fixed starting with invalid driver shows as successful ([#36](https://github.com/emccode/polly/issues/36))

### Related Bug Fixes
* Performance issue for libstorage that affected Polly ([#490](https://github.com/emccode/rexray/issues/490)), ([#218](https://github.com/emccode/libstorage/pull/218))

## Version 0.1.0 (2016/05/17)

Initial Release
