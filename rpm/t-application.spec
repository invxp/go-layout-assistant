Name: t-application
Version:0.0.1
Release: %(echo $RELEASE)%{?dist} 
Summary: application
Group: devel
License: GPL

Requires(post): chkconfig
Requires(preun): chkconfig, initscripts

AutoReqProv: none

%define _binaries_in_noarch_packages_terminate_build   0

%define _nick   		application
%define _dir			/home/application
%define _config		 	config.toml
%define _service 		%{_nick}.service

%define _systemd_dir   	/etc/systemd/system
%define _systemd_file 	init/%{_service}

BuildArch:noarch

%description
Application

%prep

%install

mkdir -p ${RPM_BUILD_ROOT}%{_dir}
mkdir -p ${RPM_BUILD_ROOT}%{_systemd_dir}

cd $OLDPWD
#cd $OLDPWD/../

bash ./build/make

%{__install} -p -m 0755 %{_systemd_file} ${RPM_BUILD_ROOT}%{_systemd_dir}/%{_service}

%clean
rm -rf ${RPM_BUILD_ROOT}

%post
systemctl stop %{_service}
systemctl enable %{_service}
systemctl start %{_service}

%preun
if [ $1 = 0 ]; then
    systemctl disable %{_service}
    systemctl stop %{_service}
fi

%files

%defattr(-,root,root)

%{_systemd_dir}/%{_service}
%{_dir}/%{_nick}
%{_dir}/resource/*

%config(noreplace) %{_dir}/%{_config}
