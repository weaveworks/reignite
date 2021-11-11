from metal import Welder


class Test:
    def __init__(self, auth_token, config):
        self.testCfg = config['test']
        devCfg = config['device']

        self.welder = Welder(auth_token, config)
        self.prj_name = config['project']
        self.key_name = devCfg['ssh_key_name']
        self.dev_name = devCfg['name']
        self.skip_teardown = self.testCfg['skip_teardown']
        self.dev_id = devCfg['id']
        self.dev_ip = None
        self.project = None
        self.key = None
        self.device = None

    def __enter__(self):
        return self

    def __exit__(self, *args, **kwargs):
        if self.skip_teardown == False:
            self.teardown()

    def setup(self):
        if self.dev_id != None:
            self.fetch_infra()
        else:
            self.create_infra()

    def run_tests(self):
        cmd = ['./test/e2e/test.sh',
               '-level.flintlockd', self.testCfg['flintlock_log_level'],
               '-level.containerd', self.testCfg['containerd_log_level'],
               ]
        if self.testCfg['skip_delete']:
            cmd.append('-skip.teardown')
            cmd.append('-skip.delete')
        if self.testCfg['skip_dmsetup']:
            cmd.append('-skip.setup.thinpool')
        try:
            self.welder.run_ssh_command(cmd, "/root/work/flintlock", False)
        except RuntimeError as e:
            print(str(e))
            pass

    def teardown(self):
        self.welder.delete_all(self.project, self.device, self.key)

    def create_infra(self):
        self.dev_ip = self.welder.create_all()

    def fetch_infra(self):
        try:
            ip = self.welder.get_device_ip(self.dev_id)
        except:
            raise
        self.ip = ip

    def device_details(self):
        return self.dev_id, self.dev_ip
