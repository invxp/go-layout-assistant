Name: t-assistant
Version:0.0.1
Release: %(echo $RELEASE)%{?dist} 
Summary: assistant
Group: application
License: GPL

Requires(post): chkconfig
Requires(preun): chkconfig, initscripts

AutoReqProv: none

%define _nick   assistant
%define _config config.toml
%define _service %{_nick}.service
%define _dir    /home/

%define _prefix %{_dir}%{_nick}
%define _init   /usr/lib/systemd/system/%{_service}

%define _service resource/%{_service}

BuildArch:noarch

%description
assistant

%prep

mkdir -p ${RPM_BUILD_ROOT}%{_prefix}
mkdir -p ${RPM_BUILD_ROOT}%{_init}

%install

cd $OLDPWD/../

bash ./build/make --prefix=${RPM_BUILD_ROOT}%{_prefix}

%{__install} -p -m 0755 %{initd} ${RPM_BUILD_ROOT}%{_init}/%{_nick}

%clean
rm -rf ${RPM_BUILD_ROOT}

%post
systemctl stop %{_service}
pkill -9 %{_nick}
systemctl enable %{_service}
systemctl start %{_service}

%preun
if [ $1 = 0 ]; then
    systemctl disable %{_service}
    systemctl stop %{_service}
fi

%files

%defattr(-,root,root)

%{_init}
%{_prefix}/%{_nick}
%{_prefix}/resource/*

%config(noreplace) %{_prefix}/%{_config}
