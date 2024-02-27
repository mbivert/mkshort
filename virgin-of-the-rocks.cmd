ffmpeg -y \
	-r 30 -t 36.80 -loop 1 -i "virgin-of-the-rocks-paris.jpg" \
	-i "/home/mb/gits/mkshort/.cache/5ab0aeb863fa30690ac886848568294f088112b3e214f1d76057b4813157ee92.png" \
	-i "/home/mb/gits/mkshort/.cache/4b2f909c63713285cbc12433a7825417f67af53f0aa4f886d8d337544af784c6.png" \
	-i "/home/mb/gits/mkshort/.cache/3f8b97412fcab49d55756a1a7499cf8b415ba6c0d9fbefc5b8692b0f4e537489.png" \
	-i "/home/mb/gits/mkshort/.cache/5c27a327e3259bef5c6b041fa599920f3fae8633ee438ddef7d3ca930166cf73.png" \
	-i "/home/mb/gits/mkshort/.cache/e3f9bb9851eb805e7ba3a1cbbd96fcf21ae9f8caba1b168b49310e6b537a96b0.png" \
	-i "/home/mb/gits/mkshort/.cache/3ff769731dbab97c486431ecf88e2b41190bec481a3ee4e1fa531d26fee365ef.png" \
	-i "/home/mb/gits/mkshort/.cache/4e882109b94eb4fe10dac69753922045dc82a329c4d0f9136538e258389858d3.png" \
	-i "/home/mb/gits/mkshort/.cache/0616f5379b27407cdb10fc28b6764b7e1f895d06a616fd788b58bae027e4c57c.png" \
	-r 30 -t 17.30 -loop 1 -i "virgin-of-the-rocks-london.jpg" \
	-i "/home/mb/gits/mkshort/.cache/10df258c045e90a9bdba6ba868ea3df937e9f180d022e23d7aed79bebb00d923.png" \
	-i "/home/mb/gits/mkshort/.cache/7999bf9bdfe9f9715ee6d1e5cff7a02f137224dcf6a86cf3ef483be89fafa4cd.png" \
	-i "/home/mb/gits/mkshort/.cache/2acc744f08511fb3ae5dfc48bb97c5fd068479a83d92d4ac1273df60262ca3c3.png" \
	-i "BMC19T1VivaldiSeasonsSpring.mp3" \
	-filter_complex "
		[0:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in0];
		[9:v] scale=1080:1920:force_original_aspect_ratio=decrease,pad=1080:1920:(ow-iw)/2:(oh-ih)/2:color=black [in1];

		[in0] [1:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.20)' [in0,1];
		[in0,1] [2:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,4.90,7.90)' [in0,2];
		[in0,2] [3:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.60,11.60)' [in0,3];
		[in0,3] [4:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,12.30,15.30)' [in0,4];
		[in0,4] [5:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,16.00,19.00)' [in0,5];
		[in0,5] [6:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,19.70,24.70)' [in0,6];
		[in0,6] [7:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,25.40,28.40)' [in0,7];
		[in0,7] [8:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,31.80,35.80)' [in0,9];
		[in1] [10:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,0.70,4.70)' [in1,1];
		[in1,1] [11:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,8.10,10.60)' [in1,3];
		[in1,3] [12:v] overlay=x=(W-w)/2:y=(H-h)/2:enable='between(t,11.30,14.30)' [in1,4];

		[in0,9] [in1,4] concat=n=2:v=1:a=0:unsafe=1 [v];
		[13:a] afade=type=in:start_time=0:duration=4.00, afade=type=out:start_time=50.10:duration=4.00 [a]
	" \
	-pix_fmt yuv420p -r 30 -movflags faststart -map "[v]" -map "[a]" -c:a aac -shortest "/dev/stdout"
