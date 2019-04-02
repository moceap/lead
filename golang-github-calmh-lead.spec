%global repo    lead
%global goipath github.com/calmh/%{repo}
%global commit 689e6190233d193c0db341d9e5bcbcfffb355738

Name:           golang-github-calmh-lead
Summary:        Lead Energy wireless LED
Version:        0
Release:        1.git%(c=%{commit}; echo ${c:0:7})%{?dist}
License:        MIT
BuildRequires:  golang-github-alecthomas-kingpin-unit-test-devel



%gometa

URL:            %{gourl}
Source0:        lead-%{commit}.zip

%description
%{summary}


%package        devel
Summary:        %{summary} devel
BuildArch:      noarch

%description    devel
%{summary}

This package contains devel files.


%prep
%forgeautosetup -p1 -n %{repo}-%{commit}


%build
%gobuildroot

%gobuild -o _bin/cmd/lead %{goipath}/cmd/lead


%install
install -d -p %{buildroot}%{_bindir}
install -p -m 0755 _bin/cmd/lead %{buildroot}%{_bindir}

%goinstall


%check
%gochecks


%files
%license LICENSE
%doc README.md

%{_bindir}/lead


%files devel -f devel.file-list
%license LICENSE
%doc README.md


%changelog
* Sun Mar 24 2019 Mosaab Alzoubi <moceap@hotmail.com> - 0-1.git689e619
- Initial
