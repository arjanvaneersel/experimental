###
### Black box testing to see if we properly emulate the CKAN API
###

import unittest
import ckanclient

class TestCkanClient(unittest.TestCase):
    def setUp(self):
        self.client = ckanclient.CkanClient("http://localhost:8080/api")

    def test_01_register_get(self):
        packages = self.client.package_register_get()
        assert isinstance(packages, list), packages
        assert len(packages) > 0

    def test_02_package_walk(self):
        for pkgid in self.client.package_register_get():
            pkg = self.client.package_entity_get(pkgid)

            assert isinstance(pkg, dict), pkg
            
            assert "id" in pkg, pkg
            assert isinstance(pkg["id"], unicode), pkg["id"]
            assert len(pkg["id"]) > 0
            
            assert "name" in pkg, pkg
            assert isinstance(pkg["name"], unicode), pkg["name"]
            assert len(pkg["name"]) > 0

            assert "revision_id" in pkg, pkg
            assert isinstance(pkg["revision_id"], unicode), pkg["revision_id"]
            assert len(pkg["revision_id"]) > 0
            
if __name__ == '__main__':
    unittest.main()
