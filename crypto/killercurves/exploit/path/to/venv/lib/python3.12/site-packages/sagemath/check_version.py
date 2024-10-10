#!/usr/bin/env sage -python

from pkg_resources import parse_version

# Here we test that the version is the correct one
def check_version(version):
    try:
        from sage.all import version as sage_version_function
    except ImportError:
        raise ValueError("Sage does not seem to be installed in this system. Please visit www.sagemath.org to fix this!")
    version = str(version)
    installed_version = sage_version_function().replace(',','').split()[2]
    if version.find('==') != -1:
        version = version.replace('==','')
        if parse_version(version) != parse_version(installed_version):
            raise ValueError("""\n******************************************************************\n
Sage version (= %s) is different from required one (= %s)\n
******************************************************************"""%(installed_version,version))
    elif version.find('>=') != -1:
        version = version.replace('>=','')
        if parse_version(version) > parse_version(installed_version):
            raise ValueError("""\n******************************************************************\n
Sage version (= %s) is older than the required one (= %s)\n
******************************************************************"""%(installed_version,version))
    elif version == '':
        pass
    else:
        raise ValueError("Version argument (=%s) not understood"%version)
    return


