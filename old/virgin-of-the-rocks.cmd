ffmpeg -y \
-loop 1 \
-i "/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-1080x1920.jpg" \
-framerate 30 \
-vf "scale=1080x1920" \
-movflags faststart \
-vf "[in]
drawtext=enable='between(t,0.7,4.2)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='Here is what seems
to be
a little-known fact',
drawtext=enable='between(t,4.9,7.9)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='About this famous
Leonardo painting,',
drawtext=enable='between(t,8.6,11.6)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='« Virgin of the Rocks »:',
drawtext=enable='between(t,12.3,15.3)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='Of the four characters',
drawtext=enable='between(t,16.0,19.0)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='Jésus, the
bottom right child',
drawtext=enable='between(t,19.7,23.7)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='Is the only one
looking at the
« Source of Light »',
drawtext=enable='between(t,24.4,27.4)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='Pay attention to his
chest and eyes',
drawtext=enable='between(t,30.8,34.8)':fontfile=/home/mb/gits/site-theme/static/fonts/cmunrm.woff:fontcolor=white:fontsize=h/20:box=1:boxcolor=black@0.5:boxborderw=10:x=(w-text_w)/2:y=(h-text_h+1)/2+t:text_align=M:text_align=C:text='The source of light can
be deduced from the
highlights and shadow
patterns.'
[out]" \
-t 36.5 \
/home/mb/gits/posts/virgin-of-the-rocks/virgin-of-the-rocks-1080x1920.mp4
