# Amazon Lightsail ���g�����A�X�}�z�̌��݈ʒu�̕\�����@

PruneMobile�́A�V�~�����[�^���Ōv�Z�����ʒu�����A�u���E�U�ŕ\�����邱�Ƃ�ړI�Ƃ������̂ł����A������A�����̃X�}�z�̈ʒu�̌��m�ɂ��g����悤�ɂ��܂���(�v����Ɂu�R�R�Z�R���v�Ƃ��Ă��g����A�Ƃ������Ƃł�)

�������������ׂɂ́A�C���^�[�l�b�g���(�N���E�h)�T�[�o��u���Ȃ���΂Ȃ�܂���BAWS��VPS(���z��p�T�[�o�[)���v�����܂����AAWS��EC2�͉^�p���ʓ|�ȏ�Ɏg�p���������ł��B�����Łu���z 500 �~�Ŏg����AWS�N���E�h��VPS�v���g�����@�ɂ��ċL�ڂ��Ă����܂��B

- Amazon Lightsail �̗����グ���@�ɂ��ẮA<a href="https://wp.kobore.net/�]�[����̋Z�p����/post-1513/" target="_blank">������</a> ���Q�l�ɂ��ĉ������B

- �����ł́A"sea-anemone.tech"�Ƃ����ˋ�̃h���C�����Ƃ��Ďg���Ă��܂����A�O��(�Ⴆ�΁u�����O.com�v)�Ńh���C���𓾂��ꍇ�́A���̖��O�ɒu�������ēǂ�ŉ������B


- ���J���̎擾���@�ɂ��ẮA<a href="https://wp.kobore.net/�]�[����̋Z�p����/post-1550/" target="_blank">������</a>���Q�l�ɂ��ĉ�����(�����ɋL�ڂ���Ă���A"go_template/server_test"�́A"PruneMobile\vps_server"�Ɠǂ݊����ĉ�����)


## Step 1 �T�[�o�̋N��

Amazon Lightsail�̃V�F������K���ȃV�F���𗧂��グ��
```
$ cd PruneMobile\vps_server    (�]�[�̊��ł́A~/go_template/server_test/ )
$ go run serverXX.go (X�͐���)
```

�ƋN�����ĉ������B

## Step 2 �n�}���(�}�[�J�\�����)�̋N��
Chromo�u���E�U(���̃u���E�U�̂��Ƃ͒m���)����A
```
https://sea-anemone.tech:8080/
```
�Ɠ��͂��ĉ������B���݂́A�����̂���n�悪�\������܂����AserverXX.go �̒��ɋL�ڂ�Ă���A�ʒu���A35.60000, 139.60000 ��Ђ��ς�����A�C�ӂ̈ʒu(����̈ʒu��)�ɕύX���邱�ƂŁA����t�߂ł̎��؎������ł��܂��B
����̏��́AGoogleMAP����擾�ł��܂��B

## Step 3 �ړ��I�u�W�F�N�g(�}�[�J�̑Ώ�)�̋N��
�X�}�z�̃u���E�U����A
```
https://sea-anemone.tech:8080/smartphone
```
�Ƃ��āA[open]�{�^���������ĉ������B�X�}�z�ňʒu���ʂ��J�n����܂�(���̍ہA�ʒu����񋟂��ėǂ����A�ƕ�������邱�Ƃ�����܂��̂ŁA"OK"�Ƃ��ĉ�����)�B
[close]�{�^������������ƒn�}��ʂ���}�[�J�������܂��B

## �����_�Ŋm�F���Ă�����_�ŁA�����꒼������

- ���[�J����js(javascript)�̃��[�f�B���O�Ɏ��s�����ׁA�]�[�̃v���C�x�[�g�T�[�o(kobore.net)���烍�[�f�B���O���Ă���BPruneMobile\vps_server\serverXX.go�̈ȉ����Q��
```
	<script src="http://kobore.net/PruneCluster.js"></script>
	<link rel="stylesheet" href="http://kobore.net/examples.css"/>
```
- ���쒆��websocket���ؒf���Ă��܂�����(�X�}�z�̕���A�ʂ̃u���E�U��ʂ𗧂��グ����)�A�I�u�W�F�N�g�����u����āA�V�X�e���S�̂������Ȃ��Ȃ�


## ���[�J���ł̒ʐM�̃y�A
server22.go �� client9.go 

## �f���V�~�����[�V�����̃y�A
server22-1.go �� pm_proxy3_1_socket.go

