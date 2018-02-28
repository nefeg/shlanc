#!/usr/bin/env bash
# project root
#   |_src
#   |_WORK_DIR
#       |_SOURCE_DIR
#       |_TMP_DIR
#       |_BUILD_NAME.tar.gz
#
#

#
# Usage ./this-script <version> <isTest>
# Example: ./this-script 0.25-0ubuntu1 0
#

# date ; sudo service ntp stop ; sudo ntpdate -s time.nist.gov ; sudo service ntp start ; date
#
# *** DIR: %project root%
#

###################################
#   checking version              #
###################################
if [[ $1 == '' ]];then
    echo "expected version (ex.: 0.25-0ubuntu1)"
    exit 1
fi


VERSION="$1"
IS_TEST="$2"
BUILD_NAME="shlanc-${VERSION}"
BUILD_NAME_ALT="shlanc_${VERSION}"
WORK_DIR=$(pwd)"/build"
TMP_DIR=$WORK_DIR/tmp
SOURCE_DIR=$WORK_DIR/source

export DEBEMAIL="evgeny.nefedkin@umbrella-web.com"
export DEBFULLNAME="Evgeny Nefedkin"


###################################
#   making workspace              #
###################################
#git archive --format=tar.gz -o $WORK_DIR/$BUILD_NAME.tar.gz ppa
## copy sources
mkdir -p $SOURCE_DIR
cp -rf src/ $SOURCE_DIR/src/
cp Makefile $SOURCE_DIR/Makefile
cp config.json $SOURCE_DIR/config.json
cp readme.md $SOURCE_DIR/readme.md

## remove .git's directories
find $SOURCE_DIR/ -name .git | rm -fr

## make archive
#
# *** DIR: %project root%/$WORK_DIR/$SOURCE_DIR
#
cd $SOURCE_DIR
tar -czf ../$BUILD_NAME.tar.gz ./*


###################################
#   make skeleton                 #
###################################
mkdir $TMP_DIR -p
cd $TMP_DIR;
dh_make -f ../$BUILD_NAME.tar.gz -s -e evgeny.nefedkin@umbrella-web.com -c gpl2 -y --createorig -p $BUILD_NAME_ALT


###################################
#   edit 'debian/control' file    #
###################################
sed -i 's/#Vcs-Browser: http:\/\/git.debian.org\/?p=collab-maint\/shlanc.git;a=summary/Vcs-Browser: https:\/\/verefkin@bitbucket.org\/verefkin\/hrentab.git/' debian/control
sed -i 's/#Vcs-Git: git:\/\/git.debian.org\/collab-maint\/shlanc.git/Vcs-Git: https:\/\/verefkin@bitbucket.org\/verefkin\/hrentab.git/' debian/control
sed -i 's/<insert the upstream URL, if relevant>/https:\/\/verefkin@bitbucket.org\/verefkin\/hrentab.git/' debian/control
#sed -ri 's/^Depends: .*/Depends: /' debian/control
sed -i 's/<insert up to 60 chars description>/Distributed and concurrency manager of deffer jobs/' debian/control
sed -i 's/<insert long description, indented with spaces>/Distributed and concurrency manager of deffer jobs and something else/' debian/control
sed -i 's/Section: unknown/Section: utils/' debian/control
sed -i 's/Build-Depends: debhelper (>= 8.0.0)/Build-Depends: debhelper (>= 9.0.0), golang (>= 1.9)/' debian/control
sed -i 's/Standards-Version: 3.9.4/Standards-Version: 3.9.5/' debian/control


###################################
#   edit 'debian/copyright' file  #
###################################
sed -i 's/<url:\/\/example.com>/https:\/\/verefkin@bitbucket.org\/verefkin\/hrentab.git/' debian/copyright
sed -i "s/<years> <put author.s name and email here>/${DEBFULLNAME} <${DEBEMAIL}>/" debian/copyright
sed -i 's/<years> <likewise for another author>/\n/' debian/copyright
sed -ri 's/^#.+//' debian/copyright

###################################
#   edit 'debian/changelog' file  #
###################################
sed -i 's/unstable;/trusty;/' debian/changelog

if [[ "$IS_TEST" != "1" ]];then
    vi debian/changelog
fi


## Remove useless files
rm debian/*.ex &>/dev/null
rm debian/*.EX &>/dev/null
#
### build package (https://help.launchpad.net/Packaging/PPA/BuildingASourcePackage)
#debuild -S -sa
#
### Upload to PPA
#dput -d ppa:onm/shlanc ${WORK_DIR}/${BUILD_NAME_ALT}-1_source.changes
#
#rm -rf $WORK_DIR

